package user

import (
	"time"

	"github.com/WeiXinao/daily_your_go/app/pkg/translator/ginx"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares"
	"github.com/WeiXinao/daily_your_go/pkg/common/core"
	timePkg "github.com/WeiXinao/daily_your_go/pkg/common/time"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}

func (uc *userServer) UpdateUser(ctx *gin.Context) {
	var form UpdateUserForm
	if err := ctx.ShouldBind(&form); err!= nil {
		ginx.HandleValidatorError(ctx, uc.translator, err)
		return
	}

	//TODO
	userID := ctx.GetFloat64(middlewares.KeyUserID)
	ud, err := uc.svc.Get(ctx, uint64(userID))
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	ud.NickName = form.Name
	loc, err := time.LoadLocation("local")
	if err!= nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrUnknown, "load location err: %w", err),nil)
		return
	}
	t, err := time.ParseInLocation(time.DateOnly, form.Birthday, loc)
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrUnknown, "日期解析错误，err: %w", err), nil)
		return
	}
	ud.Birthday = timePkg.Time{Time: t}
	ud.Gender = form.Gender
	err = uc.svc.Update(ctx, ud)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, "更新成功"), nil)
}