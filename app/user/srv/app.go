package srv

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/app/user/srv/config"
	gapp "github.com/WeiXinao/daily_your_go/gmicro/app"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_your_go/pkg/app"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

var ProviderSet = wire.NewSet(NewRegistrar, NewUserRPCServer, NewUserApp)

// controller（参数校验） -> service（具体的业务逻辑） -> data（数据库的接口）
func NewApp(name string) *app.App {
	cfg := config.New()
	return app.NewApp(
		name,	
		"daily_your_go",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
	)
}

func NewRegistrar(registry *options.RegisteryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

func NewUserApp(
	logOpts *log.Options,
	registry registry.Registrar,
	serverOpts *options.ServerOptions,
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

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		userApp, err := initApp(cfg.Log, cfg.MySQL, cfg.Telemtry, cfg.Server, cfg.Registry)
		if err != nil {
			return err
		}

		// 启动
		if err := userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}

		return nil
	}
}