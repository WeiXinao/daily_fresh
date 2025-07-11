// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package srv

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/controller/user"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/data/v1/db"
	"github.com/WeiXinao/daily_your_go/app/user/srv/internal/service/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/app"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

// Injectors from wire.go:

func initApp(logOpts *log.Options, mysqlOpts *options.MysqlOptions, telemetryOpts *options.TelemtryOptions, serverOpts *options.ServerOptions, registryOpts *options.RegisteryOptions, nacosOpts *options.NacosOptions) (*app.App, error) {
	registrar := NewRegistrar(registryOpts)
	gormDB, err := db.GetDBFactoryOr(mysqlOpts)
	if err != nil {
		return nil, err
	}
	userStore := db.NewUsers(gormDB)
	userSrv := service.NewUserService(userStore)
	userServer := user.NewUserServer(userSrv)
	nacosDataSource, err := NewNacosDatasource(nacosOpts)
	if err != nil {
		return nil, err
	}
	server, err := NewUserRPCServer(telemetryOpts, serverOpts, userServer, nacosDataSource)
	if err != nil {
		return nil, err
	}
	appApp, err := NewUserApp(logOpts, registrar, serverOpts, server)
	if err != nil {
		return nil, err
	}
	return appApp, nil
}
