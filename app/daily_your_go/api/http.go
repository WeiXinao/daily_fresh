package admin

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/config"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	aRestServer := restserver.NewServer(
		restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddewares(cfg.Server.Middlewares),
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics),
	)

	// 配置好路由
	initRouter(aRestServer, cfg)

	return aRestServer, nil
}
