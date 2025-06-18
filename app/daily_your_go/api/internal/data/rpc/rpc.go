package rpc

import (
	"sync"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/hashicorp/consul/api"
)

var _ data.DataFactory = (*grpcData)(nil)

type grpcData struct {
	uc upb.UserClient
	gc gpb.GoodsClient
}

// Goods implements data.DataFactory.
func (g *grpcData) Goods() gpb.GoodsClient {
	return g.gc
}

// User implements data.DataFactory.
func (g *grpcData) User() data.UserData {
	return NewUsers(g.uc)
}

func NewDiscovery(opts *options.RegisteryOptions) registry.Discovery {
	cfg := api.DefaultConfig()
	cfg.Address = opts.Address
	cfg.Scheme = opts.Scheme
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return consul.New(cli, consul.WithHealthCheck(true))
}

var (
	dataFactory data.DataFactory
	once sync.Once
)

// rpc 的连接，基于服务发现
func GetDataFactoryOr(options *options.RegisteryOptions) (data.DataFactory, error) {
	if options == nil && dataFactory == nil {
		return nil, errors.New("fail to get grpc data factory")
	}
	if dataFactory != nil {
		return dataFactory, nil
	}

	once.Do(func() {
		// 这里负责依赖所有的 rpc 连接
		discovery := NewDiscovery(options)
		userClient := NewUserServiceClient(discovery)
		goodsClient := NewGoodsServiceClient(discovery)
		dataFactory = &grpcData{
			uc: userClient,
			gc: goodsClient,
		} 
	})
	if dataFactory == nil {
		return nil, errors.WithCode(code.ErrGrpcConn, "fail to get grpc data factory")
	}

	return dataFactory, nil
}
