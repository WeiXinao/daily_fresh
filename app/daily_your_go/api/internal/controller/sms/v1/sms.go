package sms

import (
	"time"

	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/sms/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	baseCode "github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/app/pkg/translator/ginx"
	"github.com/WeiXinao/daily_your_go/pkg/common/core"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/storage"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type SmsController struct {
	srv sms.SmsSrv
	trans ut.Translator
}

func NewSmsController(srv sms.SmsSrv, trans ut.Translator) *SmsController {
	return &SmsController{
		srv: srv,
		trans: trans,
	}
}

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` 
	Type uint `form:"type" json:"type" binding:"required,oneof=1 2"`
	// 1. 注册发送短信验证码和动态验证码登录发送验证码
}

func (sc *SmsController) SendSms(ctx *gin.Context) {
	// 1. 参数校验
	sendSmsForm := SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		ginx.HandleValidatorError(ctx, sc.trans, err)
	}
	// 2. 生成验证码
	smsCode := sms.GenerateSmsCode(6)

	// 3. 发送验证码
	err := sc.srv.SendSms(ctx, sendSmsForm.Mobile, "SMS_181850725", "{\"code\":"+smsCode+"}")
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	// 4. redis 将验证码保存起来
	rstore := storage.RedisCluster{}
	err = rstore.SetKey(ctx, sendSmsForm.Mobile, smsCode, 5*time.Minute)
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(baseCode.ErrSuccess, "发送成功"), nil)
}