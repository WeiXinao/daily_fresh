package rpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/clientinterceptors"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/resolver/discovery"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	gis "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	timeout  time.Duration

	// discovery 接口
	discovery registry.Discovery

	unaryInterceptors  []grpc.UnaryClientInterceptor
	streamInterceptors []grpc.StreamClientInterceptor
	grpcOpts           []grpc.DialOption
	balancerName       string

	logger log.LogHelper
	enableTracing bool
	enableMetrics bool
}

// 是否开启链路追踪
func WithEnableMetrics(enableMetrics bool) ClientOption {
	return func(o *clientOptions) {
		o.enableMetrics = enableMetrics
	}
}

func WithEnableTracing(enableTracing bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enableTracing
	}
}

// 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// 设置拦截器
func WithClientUnaryInterceptors(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInterceptors = append(o.unaryInterceptors, in...)
	}
}

// 设置stream拦截器
func WithClientStreamInterceptors(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInterceptors = append(o.streamInterceptors, in...)
	}
}

// 设置 grpc 的 dial 选项
func WithDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = append(o.grpcOpts, opts...)
	}
}

// 设置负载均衡器
func WithBalacerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

// 设置日志
func WithLogger(log log.LogHelper) ClientOption {
	return func(o *clientOptions) {
		o.logger = log
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption)  (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 2 * time.Second,
		balancerName: "round_robin",
		enableTracing: true,
	}

	for _, o := range opts {
		o(&options)	
	}

	// TODO 客户端默认拦截器
	options.unaryInterceptors = append(
		options.unaryInterceptors, 
		clientinterceptors.TimeoutInterceptor(options.timeout),
	)

	// 链路追踪中间件
	if options.enableTracing {
		options.unaryInterceptors = append(
			options.unaryInterceptors, 
			otelgrpc.UnaryClientInterceptor(),
		)

		options.streamInterceptors = append(
			options.streamInterceptors, 
			otelgrpc.StreamClientInterceptor(),
		)
	}

	// metrics 中间件
	if options.enableMetrics {
		options.unaryInterceptors = append(
			options.unaryInterceptors, 
			clientinterceptors.UnaryPrometheusInterceptor,
		)
	}

	options.grpcOpts = append(
		options.grpcOpts, 
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, options.balancerName),
		),
		grpc.WithChainUnaryInterceptor(options.unaryInterceptors...),
		grpc.WithChainStreamInterceptor(options.streamInterceptors...),
	) 
	
	// TODO 服务发现的选项
	if options.discovery != nil {
			options.grpcOpts = append(
				options.grpcOpts, grpc.WithResolvers(
					discovery.NewBuilder(
						options.discovery,
						discovery.WithInsecure(insecure),
					),
				),
			)
	}

	if insecure {
		options.grpcOpts = append(
			options.grpcOpts, 
			grpc.WithTransportCredentials(gis.NewCredentials()),
		)
	}

	return grpc.DialContext(ctx, options.endpoint, options.grpcOpts...)
}