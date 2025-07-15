package server

import (
	"github.com/WeiXinao/daily_fresh/app/product/server/config"
	"github.com/WeiXinao/daily_fresh/pkg/app"
	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	gapp "github.com/WeiXinao/daily_fresh/pkg/gmicro/app"
	"github.com/WeiXinao/daily_fresh/pkg/log"
)


// controller（参数校验） -> service（具体的业务逻辑） -> data（数据库的接口）
func NewApp(name string) *app.App[*config.Config] {
	cfg := config.New()
	runFunc, productApp, subr := run(cfg)
	return app.NewApp(
		name,	
		"daily_fresh",
		app.WithOptions(cfg),
		app.WithRunFunc[*config.Config](runFunc),
		app.WithStopFunc[*config.Config](func() error {
			return productApp.Shutdown()
		}),
		app.WithSubscribeInitFunc(func(cfg *config.Config) (subscriber.Subscriber, error) {
			return subr, nil
		}),
	)
}

func run(cfg *config.Config) (app.RunFunc, *gapp.App, *subscriber.NacosSubscriber) {
	var (
		productApp *gapp.App
		subscriber *subscriber.NacosSubscriber
	)
	
	return func(basename string) error {
		appAndSubscriber, err := initApp(
			cfg.Log,
			cfg.MySQL,
			cfg.Telemtry,
			cfg.Server,
			cfg.Registry,
			cfg.Nacos,
		)
		if err != nil {
			return err
		}
		productApp = appAndSubscriber.App
		subscriber = appAndSubscriber.Subscriber

		// 启动
		if err = productApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}
		return nil
	}, productApp, subscriber
}