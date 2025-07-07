package rpc

import (
	"context"
	"time"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/clientinterceptors"
)

var _ data.UserData = (*users)(nil)

type goods struct {
	uc gpb.GoodsClient
}

const goodsServiceName = "discovery:///daily-your-go-goods-srv"

func NewGoodsServiceClient(r registry.Discovery) gpb.GoodsClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(), 
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientTimeout(1 * time.Hour),
		rpcserver.WithEndpoint(goodsServiceName),
		rpcserver.WithClientUnaryInterceptors(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}

	return gpb.NewGoodsClient(conn)
}