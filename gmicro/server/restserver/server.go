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
	"github.com/penglongli/gin-metrics/ginmetrics"
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

	// 是否开启 metrics 接口，默认开启，如果开启会自动添加 /metrics 接口
	enableMetrics bool

	// 中间件
	middleware        []string
	customMiddlewares []gin.HandlerFunc

	// jwt 信息
	jwt *JwtInfo

	// 翻译器，默认 zh
	transName string
	trans     ut.Translator

	server *http.Server

	serviceName string
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
		transName: "zh",
		Engine: gin.New(),
		serviceName: "gmicro",
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

	if len(s.customMiddlewares) > 0 {
		s.Use(s.customMiddlewares...)
	}

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

	address := fmt.Sprintf(":%d", s.port)
	// 因为 tracing 的 middleware 要使用 addr 故在此use
	s.Use(middlewares.TracingHandler(address))	

	if s.enableMetrics {
		m := ginmetrics.GetMonitor()
		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		// set middleware for gin
		m.Use(s)
	}

	_ = s.SetTrustedProxies(nil)
	s.server = &http.Server{
		Addr: address,
		Handler: s.Engine,
	}

	log.Infof("rest server is running on port: %d", s.port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("fail to start rest server")
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("rest server is stopping")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("rest server shutdown error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}