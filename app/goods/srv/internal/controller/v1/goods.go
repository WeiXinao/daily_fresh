package controller

import (
	"context"

	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/dto"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/service/v1"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/xkit/slice"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ proto.GoodsServer = (*goodServer)(nil)

type goodServer struct {
	proto.UnimplementedGoodsServer
	srv service.ServiceFactory
}

// BannerList implements proto.GoodsServer.
func (g *goodServer) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	panic("unimplemented")
}

// BatchGetGoods implements proto.GoodsServer.
func (g *goodServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goods, err := g.srv.Goods().BatchGet(ctx, slice.Map(req.GetId(), func(idx int, src int32) uint64 { return uint64(src) }))
	if err != nil {
		return nil, err
	}
	return &proto.GoodsListResponse{
		Total: int32(len(goods)),
		Data: slice.Map(goods, func(idx int, src *dto.GoodsDTO) *proto.GoodsInfoResponse {
			return ModelToResponse(src)
		}),
	}, nil
}

// BrandList implements proto.GoodsServer.
func (g *goodServer) BrandList(context.Context, *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	panic("unimplemented")
}

// CategoryBrandList implements proto.GoodsServer.
func (g *goodServer) CategoryBrandList(context.Context, *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	panic("unimplemented")
}

// CreateBanner implements proto.GoodsServer.
func (g *goodServer) CreateBanner(context.Context, *proto.BannerRequest) (*proto.BannerResponse, error) {
	panic("unimplemented")
}

// CreateBrand implements proto.GoodsServer.
func (g *goodServer) CreateBrand(context.Context, *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	panic("unimplemented")
}

// CreateCategory implements proto.GoodsServer.
func (g *goodServer) CreateCategory(context.Context, *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	panic("unimplemented")
}

// CreateCategoryBrand implements proto.GoodsServer.
func (g *goodServer) CreateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	panic("unimplemented")
}

// CreateGoods implements proto.GoodsServer.
func (g *goodServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	panic("unimplemented")
}

// DeleteBanner implements proto.GoodsServer.
func (g *goodServer) DeleteBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteBrand implements proto.GoodsServer.
func (g *goodServer) DeleteBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteCategory implements proto.GoodsServer.
func (g *goodServer) DeleteCategory(context.Context, *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteCategoryBrand implements proto.GoodsServer.
func (g *goodServer) DeleteCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteGoods implements proto.GoodsServer.
func (g *goodServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// GetAllCategorysList implements proto.GoodsServer.
func (g *goodServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	panic("unimplemented")
}

// GetCategoryBrandList implements proto.GoodsServer.
func (g *goodServer) GetCategoryBrandList(context.Context, *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	panic("unimplemented")
}

// GetGoodsDetail implements proto.GoodsServer.
func (g *goodServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	panic("unimplemented")
}

// GetSubCategory implements proto.GoodsServer.
func (g *goodServer) GetSubCategory(context.Context, *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	panic("unimplemented")
}


// GoodsList implements proto.GoodsServer.
func (g *goodServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	rsp, err := g.srv.Goods().List(
		ctx,
		metav1.ListMeta{ 
			Page: int(req.Pages),
			PageSize: int(req.PagePerNums),
		},
		req,
		[]string{},
	)
	if err != nil {
		log.Errorf("get goods list error: %v", err)
		return nil, errors.ToGrpcError(err)
	}
	return &proto.GoodsListResponse{
		Total: int32(rsp.TotalCount),
		Data: slice.Map(rsp.Items, func(idx int, src *dto.GoodsDTO) *proto.GoodsInfoResponse {
			return ModelToResponse(src)
		}),
	}, nil
}

// UpdateBanner implements proto.GoodsServer.
func (g *goodServer) UpdateBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateBrand implements proto.GoodsServer.
func (g *goodServer) UpdateBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateCategory implements proto.GoodsServer.
func (g *goodServer) UpdateCategory(context.Context, *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateCategoryBrand implements proto.GoodsServer.
func (g *goodServer) UpdateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateGoods implements proto.GoodsServer.
func (g *goodServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	panic("unimplemented")
}

func NewGoodServer(srv service.ServiceFactory) *goodServer {
	return &goodServer{srv: srv}
}

func ModelToResponse(goods *dto.GoodsDTO) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}
