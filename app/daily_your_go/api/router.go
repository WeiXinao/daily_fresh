package admin

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/user/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func initRouter(g *restserver.Server) {
	v1 := g.Group("/v1")

	uc := user.NewUserController()
	ugroup := v1.Group("/user")
	{
		ugroup.GET("/pwd_login", uc.Login)
	}
}