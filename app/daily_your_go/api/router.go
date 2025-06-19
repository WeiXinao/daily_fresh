package api

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/config"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/sms/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data/rpc"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")

	dataFactory, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}
	svcFactory := service.NewServiceFactory(dataFactory, cfg.Sms, cfg.Jwt)

	smsCtrlr := sms.NewSmsController(svcFactory, g.Translator())
	baseGroup := v1.Group("/base")
	{
		baseGroup.POST("/send_sms", smsCtrlr.SendSms)
		baseGroup.GET("/captcha", user.GetCaptcha)
	}

	uc := user.NewUserController(g.Translator(), svcFactory)
	jwtAuth, err := newJWTAuth(cfg.Jwt)
	if err != nil {
		panic(err)
	}
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/pwd_login", uc.Login)
		userGroup.POST("/register", uc.Register)
		userGroup.GET("/detail", jwtAuth.AuthFunc(), uc.GetUserDetail)
		userGroup.PUT("/update", jwtAuth.AuthFunc(), uc.UpdateUser)
	}

	goodsCtrlr := goods.NewGoodsController(svcFactory, g.Translator())
	goodsGroup := v1.Group("/goods")
	{
		goodsGroup.GET("", goodsCtrlr.List)
	}
}