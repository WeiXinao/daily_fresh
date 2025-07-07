package user

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service"
	ut "github.com/go-playground/universal-translator"
)

type userServer struct {
	translator ut.Translator
	sf service.ServiceFactory
}

func NewUserController(trans ut.Translator, sf service.ServiceFactory) *userServer {
	return &userServer{
		translator: trans,
		sf: sf,
	}
}