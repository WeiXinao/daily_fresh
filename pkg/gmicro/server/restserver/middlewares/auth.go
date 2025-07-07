package middlewares

import "github.com/gin-gonic/gin"

// 策略模式
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

type AuthOperator struct {
	strategy AuthStrategy
}

func (ao *AuthOperator) SetStrategy(strategy AuthStrategy) {
	ao.strategy = strategy
}

func (ao *AuthOperator) AuthFunc() gin.HandlerFunc {
	return ao.strategy.AuthFunc()
}