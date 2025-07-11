package trace

import (
	"sync"

	"github.com/WeiXinao/daily_fresh/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

/*
	初始化不同的 exporter 的设置
*/

const (
	KindJaeger = "jaeger"
	KindZipkin = "zipkin"
)

var (
	// set， struct 空接口体不占内存，zerobase
	agents = make(map[string]struct{})
	lock sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()
	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(o Options) error {
	var (
		sexp sdktrace.SpanExporter
		err error
		opts  = []trace.TracerProviderOption{
			trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
			trace.WithResource(resource.NewSchemaless(
				semconv.ServiceNameKey.String(o.Name),
			)),
		}
	)
	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case KindJaeger:
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))	
			if err != nil {
				return err
			}
		case KindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		opts = append(opts, trace.WithBatcher(sexp))
	}
	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] err %v", err)
	}))
	return nil
}