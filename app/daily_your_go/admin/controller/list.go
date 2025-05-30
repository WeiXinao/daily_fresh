package controller

import (
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/gin-gonic/gin"
)

func (us *userServer) List(ctx *gin.Context) {
	log.Info("GetUserList is called")
}