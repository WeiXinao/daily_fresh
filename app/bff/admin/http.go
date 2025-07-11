package admin

import (

	"github.com/WeiXinao/daily_your_go/app/user/srv/config"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func NewUserHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	urestServer := restserver.NewServer(
		restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddewares(cfg.Server.Middlewares),
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics),
	)	

	// 配置好路由
	initRouter(urestServer)

	return urestServer, nil
}