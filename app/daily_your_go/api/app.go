package admin

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/config"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	gapp "github.com/WeiXinao/daily_your_go/gmicro/app"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	"github.com/WeiXinao/daily_your_go/pkg/app"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/daily_your_go/pkg/storage"
	"github.com/hashicorp/consul/api"
)

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

func NewAPIApp(cfg *config.Config) (*gapp.App, error) {
	// 初始化 log
	log.Init(cfg.Log)
	defer log.Flush()

	// 连接 redis
	redisOpt := cfg.Redis
	redisConfig := &storage.Config{
		Host:                  redisOpt.Host,
		Port:                  redisOpt.Port,
		Addrs:                 redisOpt.Addrs,
		MasterName:            redisOpt.MasterName,
		Username:              redisOpt.Username,
		Password:              redisOpt.Password,
		Database:              cfg.Redis.Database,
		MaxIdle:               redisOpt.MaxIdle,
		MaxActive:             redisOpt.MaxActive,
		Timeout:               redisOpt.Timeout,
		EnableCluster:         redisOpt.EnableCluster,
		UseSSL:                redisOpt.UseSSL,
		SSLInsecureSkipVerify: redisOpt.SSLInsecureSkipVerify,
		EnableTracing:         redisOpt.EnableTracing,
	}
	go storage.ConnectToRedis(context.Background(), redisConfig)

	// 实例化 HTTP 服务
	restServer, err := NewAPIHTTPServer(cfg)
	if err != nil {
		return nil, err
	}

	// 实例化服务注册中心
	registry := NewRegistrar(cfg.Registry)

	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRestServer(restServer),
		gapp.WithRegistrar(registry),
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		apiApp, err := NewAPIApp(cfg)
		if err != nil {
			return err
		}

		// 启动
		if err = apiApp.Run(); err != nil {
			log.Errorf("run api app error: %s", err)
			return err
		}

		return nil
	}
}
