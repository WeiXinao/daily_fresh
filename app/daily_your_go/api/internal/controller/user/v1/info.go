package user

import (
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares"
	"github.com/WeiXinao/daily_your_go/pkg/common/core"
	"github.com/gin-gonic/gin"
)

func (us *userServer) GetUserDetail(ctx *gin.Context) {
	userID, _ := ctx.Get(middlewares.KeyUserID)
	ud, err := us.svc.Get(ctx, uint64(userID.(float64)))
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, gin.H{
		"name": ud.NickName,
		"birthday": ud.Birthday.Format(time.DateOnly),
		"gender": ud.Gender,
		"mobile": ud.Mobile,
	})
}