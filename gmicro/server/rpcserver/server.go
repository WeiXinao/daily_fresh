package rpcserver

import (
	"net"
	"net/url"
	"time"

	apimd "github.com/WeiXinao/daily_your_go/gmicro/api/metadata"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/serverinterceptors"
	"github.com/WeiXinao/daily_your_go/pkg/host"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOption func(s *Server)

type Server struct {
	*grpc.Server

	address            string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	grpcOpts           []grpc.ServerOption
	lis                net.Listener
	timeout            time.Duration

	health   *health.Server
	metadata apimd.MetadataServer
	endpoint *url.URL
}

func (s *Server) Address() string {
	return s.address 
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address:            ":0",
		health:             health.NewServer(),
		// timeout:            1 * time.Second,
		unaryInterceptors:  []grpc.UnaryServerInterceptor{},
		streamInterceptors: []grpc.StreamServerInterceptor{},
		grpcOpts:           []grpc.ServerOption{},
	}

	for _, opt := range opts {
		opt(srv)
	}

	// TODO 我们现在希望用户不设置拦截器的情况下，我们都会默认加上一些必须的拦截器，crash，tracing
	srv.unaryInterceptors = append(
		srv.unaryInterceptors,
		serverinterceptors.UnaryCrashInterceptor,
	)
	
	srv.streamInterceptors = append(
		srv.streamInterceptors,
		serverinterceptors.StreamCrashInterceptor,
	)

	if srv.timeout > 0 {
		srv.unaryInterceptors = append(
			srv.unaryInterceptors, 
			serverinterceptors.UnaryTimeoutInterceptor(srv.timeout),
		)
	}

	// 把我们传入的拦截器转换成 grpc 的 ServerOption
	srv.grpcOpts = append(
		srv.grpcOpts,
		grpc.ChainUnaryInterceptor(srv.unaryInterceptors...),
		grpc.ChainStreamInterceptor(srv.streamInterceptors...),
	)

	srv.Server = grpc.NewServer(srv.grpcOpts...)

	// 实例化 metadata 的 Server
	if srv.metadata == nil {
		srv.metadata = apimd.NewServer(srv.Server)
	}

	// 解析 address
	err := srv.listenAndEndpoint()
	if err != nil {
		panic(err)
	}

	// 注册 health
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	// 注册 metadata 的 server
	// 可以支持用户直接通过 grpc 的一个接口查看当前支持的所有的 rpc 服务
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)

	return srv
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = append(s.unaryInterceptors, in...)
	}
}

func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = append(s.streamInterceptors, in...)
	}
}

func WithOptions(in ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, in...)
	}
}

// 完成 ip 和端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}

	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

// 启动 grpc 的服务
func (s *Server) Start() error {
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

func (s *Server) Stop() error {
	// 设置服务的状态为 not_serving，防止接收新的请求过来
	s.health.Shutdown()
	s.GracefulStop()
	log.Info("[grpc] server stopped")
	return nil
}
