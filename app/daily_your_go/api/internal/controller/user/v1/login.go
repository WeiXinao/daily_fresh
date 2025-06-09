package user

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/translator/ginx"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/gin-gonic/gin"
)

type PasswordLoginForm struct {
	// 手机号码格式有规范可循，自定义 validator
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` 
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id"	binding:"required"`
}

func (us *userServer) Login(ctx *gin.Context) {
	log.Info("Login is called")
	
	// 表单验证
	var form PasswordLoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ginx.HandleValidatorError(ctx, us.translator, err) 
	}

	// 密码验证
}