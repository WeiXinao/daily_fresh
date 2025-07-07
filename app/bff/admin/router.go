package admin

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/admin/controller"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func initRouter(g *restserver.Server) {
	v1 := g.Group("/v1")

	uc := controller.NewUserController()
	ugroup := v1.Group("/user")
	{
		ugroup.GET("/list", uc.List)
	}
}