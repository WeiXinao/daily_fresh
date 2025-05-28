package rpcserver

import (
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"google.golang.org/grpc"
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