package srv

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/app/user/server/config"
	"github.com/WeiXinao/daily_fresh/pkg/app"
	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	gapp "github.com/WeiXinao/daily_fresh/pkg/gmicro/app"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry/consul"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

var ProviderSet = wire.NewSet(NewNacosDatasource, NewRegistrar, NewUserRPCServer, NewUserApp)

// controller（参数校验） -> service（具体的业务逻辑） -> data（数据库的接口）
func NewApp(name string) *app.App[*config.Config] {
	cfg := config.New()
	var (
		userApp *gapp.App
		err error	
	)
	return app.NewApp(
		name,	
		"daily_your_go",
		app.WithOptions(cfg),
		app.WithRunFunc[*config.Config](func(basename string) error {
			userApp, err = initApp(cfg.Log, cfg.MySQL, cfg.Telemtry, cfg.Server, cfg.Registry, cfg.Nacos)
			if err != nil {
				return err
			}
			// 启动
			if err := userApp.Run(); err != nil {
				log.Errorf("run user app error: %s", err)
				return err
			}
			return nil
		}),
		app.WithStopFunc[*config.Config](func() error {
			return userApp.Shutdown()
		}),
		app.WithSubscribeInitFunc(func(cfg *config.Config) (subscriber.Subscriber, error) {
			nacos := cfg.Nacos
			sc := []constant.ServerConfig{
				{
					IpAddr: nacos.Host,
					Port: uint64(nacos.Port),
				},
			}
			cc := constant.ClientConfig{
				NamespaceId:         nacos.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "tmp/nacos/log",
				CacheDir:            "tmp/nacos/cache",
				LogLevel:            "debug",
			}
			configClient, err := clients.CreateConfigClient(map[string]interface{}{
				"serverConfigs": sc,
				"clientConfig":  cc,
			})
			if err != nil {
				return nil, err
			}
			return subscriber.NewNacosSubscriber(configClient, nacos.Group, nacos.DataId), nil
		}),
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
		userApp, err := initApp(cfg.Log, cfg.MySQL, cfg.Telemtry, cfg.Server, cfg.Registry, cfg.Nacos)
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