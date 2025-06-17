package service

import (
	"context"
	"sync"

	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	search "github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/dto"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/xkit/slice"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodsSvc interface {
	// 商品列表
	List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)

	// 商品详情
	Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error)

	// 创建商品
	Create(ctx context.Context, goods *dto.GoodsDTO) error

	// 更新商品
	Update(ctx context.Context, goods *dto.GoodsDTO) error

	// 删除商品
	Delete(ctx context.Context, id uint64) error

	// 批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

var _ GoodsSvc = (*goodsService)(nil)

type goodsService struct {
	data         data.GoodsStore
	categoryData data.CategoryStore
	searchData   search.GoodsStore
	brandData    data.BrandsStore
}

func NewGoodsService(data data.GoodsStore, categoryData data.CategoryStore, searchData search.GoodsStore, brandData data.BrandsStore) GoodsSvc {
	return &goodsService{
		data:         data,
		categoryData: categoryData,
		searchData:   searchData,
		brandData:    brandData,
	}
}

// BatchGet implements GoodsSvc.
func (g *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	// go-zero 非常好用，但是我们自己去做并发的话 - 一次性启动多个 goroutine
	var (
		ret []*dto.GoodsDTO
		mu sync.Mutex
	)
	err := mr.Finish(slice.Map(ids, func(idx int, src uint64) func() error {
		return func() error {
			goodDTO, err := g.Get(ctx, src)
			if err != nil {
				return err
			}
			mu.Lock()
			ret = append(ret, goodDTO)
			mu.Unlock()
			return nil
		}
	})...)
	if err != nil {
		return nil, err
	}
	return ret, nil
	// ds, err := g.data.ListByID(ctx, ids, []string{})
	// if err != nil {
	// 	return nil, err
	// }
	// return slice.Map(ds.Items, func(idx int, src *do.GoodsDO) *dto.GoodsDTO {
	// 	return &dto.GoodsDTO{
	// 		GoodsDO: *src,
	// 	}
	// }), nil
}

// Create implements GoodsSvc.
func (g *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	/*
		数据写 mysql，然后写 es
	*/
	_, err := g.brandData.Get(ctx, uint64(goods.BrandsID))
	if err != nil {
		return err
	}

	_, err = g.categoryData.Get(ctx, uint64(goods.CategoryID))
	if err != nil {
		return err
	}

	// 之前的入 es 的方案是给 gorm 添加 aftercreate
	// 分布式事务，异构数据库的事务，基于可靠消息最终一致性\
	// 比较重的方案：每次都要发送一个事务消息
	txn := g.data.Begin() // 非常小心
	defer func() {        // 很重要
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()
	err = g.data.CreateInTxn(ctx, txn, &goods.GoodsDO)
	if err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}

	goodsSearch := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		ShopPrice:   goods.ShopPrice,
		GoodsBrief:  goods.GoodsBrief,
	}
	err = g.searchData.Create(ctx, &goodsSearch) // 这个接口如果超时了
	if err != nil {
		txn.Rollback()
		return err
	}
	txn.Commit()
	return nil
}

// Delete implements GoodsSvc.
func (g *goodsService) Delete(ctx context.Context, id uint64) error {
	panic("unimplemented")
}

// Get implements GoodsSvc.
func (g *goodsService) Get(ctx context.Context, id uint64) (*dto.GoodsDTO, error) {
	goods, err := g.data.Get(ctx, id)
	if err != nil {
		log.Errorf("data.Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *goods,
	}, nil
}

func retrieveIDs(category *do.CategoryDO) []uint64 {
	ids := make([]uint64, 0)
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint64(category.ID))
	for _, sub := range category.SubCategory {
		ids = append(ids, retrieveIDs(sub)...)
	}
	return ids
}

// List implements GoodsSvc.
func (g *goodsService) List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	searchReq := search.GoodsFilterRequest{GoodsFilterRequest: req}
	if req.TopCategory > 0 {
		category, err := g.categoryData.Get(ctx, uint64(searchReq.TopCategory))
		if err != nil {
			log.Errorf("categoryData.Get err: %v", err)
			return nil, err
		}
		ids := retrieveIDs(category)
		searchReq.CategoryIDs = slice.Map(ids, func(idx int, src uint64) any { return src })
	}

	goodsList, err := g.searchData.Search(ctx, &searchReq)
	if err != nil {
		log.Errorf("searchData.Search err: %v", err)
		return nil, err
	}

	goodsIds := slice.Map(goodsList.Items, func(idx int, src *do.GoodsSearchDO) uint64 {
		return uint64(src.ID)
	})
	goods, err := g.data.ListByID(ctx, goodsIds, orderby)
	if err != nil {
		return nil, err
	}

	return &dto.GoodsDTOList{
		TotalCount: int(goodsList.TotalCount),
		Items: slice.Map(goods.Items, func(idx int, src *do.GoodsDO) *dto.GoodsDTO {
			return &dto.GoodsDTO{
				GoodsDO: *src,
			}
		}),
	}, nil
}

// Update implements GoodsSvc.
func (g *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	panic("unimplemented")
}
