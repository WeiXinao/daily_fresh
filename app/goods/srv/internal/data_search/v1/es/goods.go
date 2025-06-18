package es

import (
	"context"
	"encoding/json"
	"strconv"

	search "github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/xkit/slice"
	"github.com/olivere/elastic/v7"
)

var _ search.GoodsStore = (*goods)(nil)

type goods struct {
	esClient *elastic.Client
}

// Create implements search.GoodsStore.
func (g *goods) Create(ctx context.Context, goods *do.GoodsSearchDO) error {
	_, err := g.esClient.Index().
		Index(goods.GetIndexName()).
		Id(strconv.Itoa(int(goods.ID))).
		BodyJson(&goods).
		Do(ctx)
	return errors.WithCode(code.ErrElasticsearch, err.Error())
}

// Delete implements search.GoodsStore.
func (g *goods) Delete(ctx context.Context, id uint64) error {
	_, err := g.esClient.Delete().
		Index(do.GoodsSearchDO{}.GetIndexName()).
		Id(strconv.Itoa(int(id))).
		Refresh("true").
		Do(ctx)
	return errors.WithCode(code.ErrElasticsearch, err.Error())
}


// Search implements search.GoodsStore.
func (g *goods) Search(ctx context.Context, req *search.GoodsFilterRequest) (*do.GoodsSearchDOList, error) {
	// match bool 复合查询
	q := elastic.NewBoolQuery()
	// localDB := global.DB.Model(&model.Goods{})
	if req.KeyWords != "" {
		// 搜索
		// localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(req.GetKeyWords(), "name", "goods_brief"))
	}
	if req.IsHot {
		// localDB = localDB.Where("is_hot = true")
		q = q.Filter(elastic.NewTermQuery("is_hot", req.GetIsHot()))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.GetIsNew()))
	}
	if req.PriceMin > 0 {
		// localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.GetPriceMin()))
	}
	if req.PriceMax > 0 {
		// localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.GetPriceMax()))
	}

	if req.Brand > 0 {
		// localDB = localDB.Where("brands_id = ?", req.Brand)
		q = q.Filter(elastic.NewTermQuery("brand_id", req.GetBrand()))
	}

	if req.TopCategory > 0 {
		// 生成 terms 查询
		q = q.Filter(elastic.NewTermsQuery("category_id", req.CategoryIDs...))
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	} 

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums < 0:
		req.PagePerNums = 10
	}

	// 排序
	res, err := g.esClient.Search().
		Index(do.GoodsSearchDO{}.GetIndexName()).
		Query(q).
		From(int(req.Pages - 1  * req.PagePerNums)).
		Size(int(req.PagePerNums)).
		Do(ctx)
	if err != nil {
		return nil, errors.WithCode(code.ErrElasticsearch, err.Error()) 
	}

	var ret do.GoodsSearchDOList
	ret.TotalCount = res.Hits.TotalHits.Value
	ret.Items = slice.Map(res.Hits.Hits, func(idx int, src *elastic.SearchHit) *do.GoodsSearchDO {
		goods := do.GoodsSearchDO{}
		er := json.Unmarshal(src.Source, &goods)
		if er != nil {
			err = errors.WithCode(code.ErrEsUnmarshal, er.Error())
			return nil
		}
		return &goods
	})
	return &ret, err

}

// Update implements search.GoodsStore.
func (g *goods) Update(ctx context.Context, goods *do.GoodsSearchDO) error {
	err := g.Delete(ctx, uint64(goods.ID))
	if err != nil {
		return err
	}
	err = g.Create(ctx, goods)
	if err != nil {
		return errors.WithCode(code.ErrEsUnmarshal, err.Error())
	}
	return nil
}

// func NewGoods(esClient *elastic.Client) *goods {
// 	return &goods{
// 		esClient: esClient,
// 	}
// }

func newGoods(dataSearch *dataSearch) *goods {
	return &goods{
		esClient: dataSearch.esClient,
	}
}