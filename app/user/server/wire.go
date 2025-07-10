//go:build wireinject
// +build wireinject

package srv

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/app/user/server/internal/controller/user"
	"github.com/WeiXinao/daily_fresh/app/user/server/internal/data/v1/db"
	"github.com/WeiXinao/daily_fresh/app/user/server/internal/service/v1"
	gapp "github.com/WeiXinao/daily_fresh/pkg/gmicro/app"
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
) (*gapp.App, error) {
	panic(wire.Build(
		ProviderSet,
		user.ProviderSet,	
		service.ProviderSet,
		db.ProviderSet,
	))
}
