package rpc

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	"github.com/hashicorp/consul/api"
)

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

// rpc 的连接，基于服务发现
func GetDataFactoryOr(options *options.RegisteryOptions) (data.UserData, error) {
	// 这里负责依赖所有的 rpc 连接
	discovery := NewDiscovery(options)
	userClient := NewUserServiceClient(discovery)
	return NewUsers(userClient), nil
}