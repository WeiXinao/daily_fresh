package user

import (
	"net/http"

	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		log.Errorf("生成验证码错误:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath": b64s,
	})
}