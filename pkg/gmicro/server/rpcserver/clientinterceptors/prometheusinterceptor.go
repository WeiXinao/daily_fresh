package clientinterceptors

import (
	"context"
	"time"

	metric "github.com/WeiXinao/daily_your_go/gmicro/core/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/*
两个基本指标
	1. 每个请求的耗时（histogram）
	2. 每个请求的转态计数器（counter）

/user 状态码 有 label 主要是状态码
*/

const serverNamespace = "rpc_client"

var (
	metricsServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name: "daily_your_go_duration_ms",
		Help: "rpc server requests duration(ms).",
		Labels: []string{"method"},
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricsServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name: "daily_your_go_code_total",
		Help: "rpc server requests code count.",
		Labels: []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error  {
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	dur := time.Since(startTime).Milliseconds()
	code := status.Code(err).String()

	// 记录了耗时
	metricsServerReqDur.Observe(dur, method)
	
	// 记录了状态码
	metricsServerReqCodeTotal.Inc(method, code)
	return err
}