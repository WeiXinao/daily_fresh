package admin

import (
	"github.com/WeiXinao/daily_fresh/app/bff/admin/config"
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/app"
	gapp "github.com/WeiXinao/daily_fresh/pkg/gmicro/app"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry/consul"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/hashicorp/consul/api"
)

// controller（参数校验） -> service（具体的业务逻辑） -> data（数据库的接口）
func NewApp(name string) *app.App[*config.Config] {
	cfg := config.New()
	return app.NewApp(
		name,	
		"daily_your_go",
		app.WithOptions(cfg),
		app.WithRunFunc[*config.Config](run(cfg)),
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

func NewUserApp(cfg *config.Config) (*gapp.App, error) {
	// 初始化 log
	log.Init(cfg.Log)
	defer log.Flush()

	// 实例化服务
	restServer, err := NewUserHTTPServer(cfg)
	if err != nil {
		return nil, err
	}

	// 服务注册
	registry := NewRegistrar(cfg.Registry)

	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRestServer(restServer),
		gapp.WithRegistrar(registry),
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		userApp, err := NewUserApp(cfg)
		if err != nil {
			return err
		}

		// 启动
		if err = userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}

		return nil
	}
}