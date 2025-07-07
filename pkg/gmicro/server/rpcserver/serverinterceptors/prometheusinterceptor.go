package serverinterceptors

import (
	"context"
	"time"

	metric "github.com/WeiXinao/daily_fresh/pkg/gmicro/core/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/*
两个基本指标
	1. 每个请求的耗时（histogram）
	2. 每个请求的转态计数器（counter）

/user 状态码 有 label 主要是状态码
*/

const serverNamespace = "rpc_server"

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

func UnaryPrometheusInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	dur := time.Since(startTime).Milliseconds()
	code := status.Code(err).String()

	// 记录了耗时
	metricsServerReqDur.Observe(dur, info.FullMethod)
	
	// 记录了状态码
	metricsServerReqCodeTotal.Inc(info.FullMethod, code)
	return resp, err
}