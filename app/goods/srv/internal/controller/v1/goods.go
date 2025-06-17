package controller

import (
	"context"
	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ proto.GoodsServer = (*goodServer)(nil)

type goodServer struct {
	proto.UnimplementedGoodsServer
	srv service.GoodsSvc
}

// BannerList implements proto.GoodsServer.
func (g *goodServer) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	panic("unimplemented")
}

// BatchGetGoods implements proto.GoodsServer.
func (g *goodServer) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	panic("unimplemented")
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
func (g *goodServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	panic("unimplemented")
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

func NewGoodServer(srv service.GoodsSvc) *goodServer {
	return &goodServer{srv: srv}
}
