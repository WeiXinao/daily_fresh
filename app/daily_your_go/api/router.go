package admin

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/config"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/sms/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/controller/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data/rpc"
	usersvc "github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/user/v1"
	smssvc "github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/sms/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")

	smsSvc := smssvc.NewSmsService(cfg.Sms)
	smsCtrlr := sms.NewSmsController(smsSvc, g.Translator())
	baseGroup := v1.Group("/base")
	{
		baseGroup.GET("/send_sms", smsCtrlr.SendSms)
		baseGroup.POST("/captcha", user.GetCaptcha)
	}

	userData, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}
	us := usersvc.NewUserService(userData, cfg.Jwt)
	uc := user.NewUserController(g.Translator(), us)
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/pwd_login", uc.Login)
	}
}