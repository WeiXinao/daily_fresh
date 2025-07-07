package goods

import (
	"context"

	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
)

type GoodsSrv interface {
	List(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error)
}

type goodsService struct {
	data data.DataFactory
}

// List implements GoodsSrv.
func (g *goodsService) List(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	return g.data.Goods().GoodsList(ctx, request)
}

func NewGoods(data data.DataFactory) *goodsService {
	return &goodsService{data: data}
}

var _ GoodsSrv = (*goodsService)(nil)
