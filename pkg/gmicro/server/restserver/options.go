package restserver

import (
	"github.com/gin-gonic/gin"
)

type ServerOption func (s *Server)

func WithServiceName(srvName string) ServerOption {
	return func(s *Server) {
		s.serviceName = srvName
	}
}

func WithEnableMetrics(enableMetrics bool) ServerOption {
	return func(o *Server) {
		o.enableMetrics = enableMetrics
	}
}

func WithEnableProfiling(profiling bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = profiling
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.mode = mode
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMiddewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middleware = append(s.middleware, middlewares...)
	}
}

func WithCustomMiddlewares(cmws []gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		s.customMiddlewares = append(s.customMiddlewares, cmws...)
	}	
}

func WithHealthz(healthz bool) ServerOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithJwt(jwt *JwtInfo) ServerOption {
	return func(s *Server) {
		s.jwt = jwt
	}
}

func WithTransName(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}
