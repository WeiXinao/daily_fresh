package srv

import (
	"github.com/WeiXinao/daily_your_go/app/user/srv/config"
	gapp "github.com/WeiXinao/daily_your_go/gmicro/app"
	"github.com/WeiXinao/daily_your_go/pkg/app"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

// controller（参数校验） -> service（具体的业务逻辑） -> data（数据库的接口）
func NewApp(name string) *app.App {
	cfg := config.New()
	return app.NewApp(
		name,	
		"daily_your_go",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		app.WithNoConfig(),
	)
}

func NewUserApp(cfg *config.Config) (*gapp.App, error) {
	// 初始化 log
	log.Init(cfg.Log)
	defer log.Flush()

	// 实例化服务
	rpcServer, err := NewUserRPCServer(cfg)
	if err != nil {
		panic(err)
	}
	return gapp.New(gapp.WithRPCServer(rpcServer)), nil
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