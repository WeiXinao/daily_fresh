package admin

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data/rpc"
	usersvc "github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/config"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")

	userData, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}
	us := usersvc.NewUserService(userData, cfg.Jwt)
	uc := user.NewUserController(g.Translator(), us)

	baseGroup := v1.Group("/base")
	{
		baseGroup.POST("/captcha", user.GetCaptcha)
	}

	userGroup := v1.Group("/user")
	{
		userGroup.POST("/pwd_login", uc.Login)
	}
}