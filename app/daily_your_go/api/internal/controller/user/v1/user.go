package user

import ut "github.com/go-playground/universal-translator"

type userServer struct {
	translator ut.Translator
}

func NewUserController(trans ut.Translator) *userServer {
	return &userServer{
		translator: trans,
	}
}