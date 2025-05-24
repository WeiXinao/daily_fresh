package main

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

func A() error {
	return errors.WithCode(code.ErrUserNotFound, "user not found")
}

func main() {
	log.Init(log.NewOptions())
	log.Infof("hello %s", "xiaoxin")
	// 加了一些方法

	e := A()
	errors.ParseCoder(e)
}
