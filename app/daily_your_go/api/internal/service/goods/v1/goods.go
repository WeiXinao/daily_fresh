package goods

import (
	"context"

	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
)

type GoodsSrv interface {
	List(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error)
}