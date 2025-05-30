package restserver

import (
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/gin-gonic/gin"
)

type JwtInfo struct {
	// default to "JWT"
	Realm string
	// default to empty
	Key string	
	// default to 7 days
	Timeout time.Duration
	// default to 7 days
	MaxRefresh time.Duration
}

// wrapper for gin.Engine
type Server struct {
	*gin.Engine

	// 端口号，默认值 8080
	port int

	// 开发模式，默认值 debug
	mode string

	// 是否开启健康检查接口，默认开启，如果开启会自动添加 /health 接口
	healthz bool

	// 是否开启 pprof 接口，默认开启，如果开启会自动添加 /debug/pprof 接口
	enableProfiling  bool

	// 中间件
	middleware []string
	customMiddlewares []gin.HandlerFunc

	// jwt 信息
	jwt *JwtInfo
} 

// debug 模式和 release 模式的区别主要是打印的日志不同
// 环境变量的模式，在 docker，k8s部署中很常用

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		port: 8080,	
		mode: "debug",
		healthz: true,
		enableProfiling: true,
		jwt: &JwtInfo{
			"JWT",
			"xSLNbMbkh8tHu5m^pxPwwyJd!Fvrjf8G",
			7 * 24 * time.Hour,
			7 * 24 * time.Hour,
		},
	}

	for _, opt := range opts {
		opt(s)	
	}

	for _, m := range s.middleware {
		mw, ok := middlewares.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}

	s.Use(s.customMiddlewares...)
	
	return s
}