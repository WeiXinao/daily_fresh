package user

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/translator/ginx"
	"github.com/WeiXinao/daily_your_go/pkg/common/core"
	"github.com/gin-gonic/gin"
)

type RegisterForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` 
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

func (us *userServer) Register(ctx *gin.Context) {
	regForm := RegisterForm{}
	if err := ctx.ShouldBind(&regForm); err!= nil {
		ginx.HandleValidatorError(ctx, us.translator, err)
		return
	}

	ud, err := us.sf.Users().Register(ctx, regForm.Mobile, regForm.Password, regForm.Code)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, gin.H{
		"id": ud.ID,
		"nick_name": ud.NickName,
		"token": ud.Token,
		"expired_at": ud.ExpiresAt,
	})
}