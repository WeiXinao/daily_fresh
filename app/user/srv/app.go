package user

import (
	"github.com/WeiXinao/daily_your_go/app/user/srv/config"
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

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		log.Infof("%s start", basename)
		log.Info(cfg.Log.Level)
		return nil
	}
}