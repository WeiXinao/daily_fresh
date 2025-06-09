package user

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/user/v1"
	ut "github.com/go-playground/universal-translator"
)

type userServer struct {
	translator ut.Translator
	svc user.UserService
}

func NewUserController(trans ut.Translator, svc user.UserService) *userServer {
	return &userServer{
		translator: trans,
		svc: svc,
	}
}