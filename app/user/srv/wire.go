//go:build wireinject
// +build wireinject

package srv

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/controller/user"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/data/v1/db"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/service/v1"
	gapp "github.com/WeiXinao/daily_your_go/gmicro/app"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/google/wire"
)

func initApp(
	logOpts *log.Options,
	mysqlOpts *options.MysqlOptions,
	telemetryOpts *options.TelemtryOptions,
	serverOpts *options.ServerOptions,
	registryOpts *options.RegisteryOptions,
) (*gapp.App, error) {
	panic(wire.Build(
		ProviderSet,
		user.ProviderSet,	
		service.ProviderSet,
		db.ProviderSet,
	))
}
