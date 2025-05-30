package restserver

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/pprof"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/validation"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
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
	enableProfiling bool

	// 中间件
	middleware        []string
	customMiddlewares []gin.HandlerFunc

	// jwt 信息
	jwt *JwtInfo

	// 翻译器
	transName string
	trans     ut.Translator
}

// debug 模式和 release 模式的区别主要是打印的日志不同
// 环境变量的模式，在 docker，k8s部署中很常用

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		port:            8080,
		mode:            "debug",
		healthz:         true,
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

// start rest server
func (s *Server) Start(ctx context.Context) error {
	modes := []string{gin.DebugMode, gin.TestMode, gin.ReleaseMode}
	if !slices.Contains(modes, s.mode) {
		return errors.New("mode must be one of debug/release/test")
	}

	// 设置开发模式，打印路由信息
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// TODO: 初始化翻译器
	if err := s.initTrans(s.transName); err != nil {
		log.Errorf("initTrans error %s", err.Error())
		return err
	}

	// 注册 mobile 验证器
	validation.RegisterMobile(s.trans)

	// 根据配置初始化 pprof 路由
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	log.Infof("rest server is running on port: %d", s.port)
	_ = s.SetTrustedProxies(nil)
	err := s.Run(fmt.Sprintf(":%d", s.port))
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("fail to start rest server")
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	log.Infof("rest server is stopping")
	
}