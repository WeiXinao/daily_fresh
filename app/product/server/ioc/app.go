package ioc

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	gapp "github.com/WeiXinao/daily_fresh/pkg/gmicro/app"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(wire.Struct(new(AppAndSubscriber), "*"), NewProductApp)

type AppAndSubscriber struct {
	App *gapp.App
	Subscriber *subscriber.NacosSubscriber
}

func NewProductApp(
	logOpts *log.Options,
	serverOpts *options.ServerOptions,
	registry registry.Registrar,
	rpcServer *rpcserver.Server,
) (*gapp.App, error) {
	// 初始化 log
	log.Init(logOpts)
	defer log.Flush()

	return gapp.New(
		gapp.WithName(serverOpts.Name),
		gapp.WithRPCServer(rpcServer),
		gapp.WithRegistrar(registry),
	), nil
}
