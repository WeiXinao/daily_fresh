package db

import (
	"context"
	"sync"
	"time"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	ipb "github.com/WeiXinao/daily_your_go/api/inventory/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	rpc "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	_ "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/resolver/direct"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/selector"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/selector/random"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const (
	goodServiceName      = "discovery:///daily-your-go-goods-srv"
	inventoryServiceName = "discovery:///daily-your-go-inventory-srv"
)

var (
	goodsClient   gpb.GoodsClient
	invClient     ipb.InventoryClient
	registry      *consul.Registry
	goodsRpcOnce  sync.Once
	invRpcOnce    sync.Once
	registeryOnce sync.Once
)

func initServiceRegestry(registeryOpts *options.RegisteryOptions) (*consul.Registry, error) {
	if registeryOpts == nil && registry == nil {
		return nil, errors.New("init service failed")
	}
	if registry != nil {
		return registry, nil
	}
	var (
		err error
		cli *api.Client
	)
	registeryOnce.Do(func() {
		selector.SetGlobalSelector(random.NewBuilder())
		rpc.InitBuilder()

		cfg := api.DefaultConfig()
		cfg.Address = registeryOpts.Address
		cfg.Scheme = registeryOpts.Scheme
		cli, err = api.NewClient(cfg)
		registry = consul.New(cli, consul.WithHealthCheck(true))
	})
	if err != nil {
		return nil, err
	}
	return registry, nil
}

func connGoodsRpcServer(ctx context.Context, timeOut time.Duration,
	registery *consul.Registry) (gpb.GoodsClient, error) {
	if registery == nil && goodsClient == nil {
		return nil, errors.New("init goods client failed")
	}
	if goodsClient != nil {
		return goodsClient, nil
	}
	var (
		err  error
		conn *grpc.ClientConn
	)
	goodsRpcOnce.Do(func() {
		conn, err = rpc.DialInsecure(
			ctx,
			rpc.WithBalacerName("selector"),
			rpc.WithDiscovery(registery),
			rpc.WithClientTimeout(timeOut),
			rpc.WithEndpoint(goodServiceName),
		)
		goodsClient = gpb.NewGoodsClient(conn)
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return goodsClient, nil
}

func connInventoryRpcServer(ctx context.Context, timeOut time.Duration,
	registery *consul.Registry) (ipb.InventoryClient, error) {
	if invClient == nil && registery == nil {
		return nil, errors.New("init inventory client failed")
	}

	if invClient != nil {
		return invClient, nil
	}
	var (
		conn *grpc.ClientConn
		err  error
	)
	invRpcOnce.Do(func() {
		conn, err = rpc.DialInsecure(
			ctx,
			rpc.WithBalacerName("selector"),
			rpc.WithDiscovery(registery),
			rpc.WithClientTimeout(timeOut),
			rpc.WithEndpoint(inventoryServiceName),
		)
		invClient = ipb.NewInventoryClient(conn)
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return invClient, nil
}
