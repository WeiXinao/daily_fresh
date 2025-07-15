//go:build wireinject
// +build wireinject

package server

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/app/product/server/ioc"
	"github.com/WeiXinao/daily_fresh/app/product/server/ioc/handler"
	"github.com/WeiXinao/daily_fresh/app/product/server/ioc/infra"
	"github.com/WeiXinao/daily_fresh/app/product/server/ioc/repository"
	"github.com/WeiXinao/daily_fresh/app/product/server/ioc/service"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/google/wire"
)

func initApp(
	logOpts *log.Options,
	mysqlOpts *options.MysqlOptions,
	telemetryOpts *options.TelemtryOptions,
	serverOpts *options.ServerOptions,
	registryOpts *options.RegisteryOptions,
	nacosOpts *options.NacosOptions,
) (*ioc.AppAndSubscriber, error) {
	panic(wire.Build(
		ioc.ProviderSet,
		infra.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
	))
}
