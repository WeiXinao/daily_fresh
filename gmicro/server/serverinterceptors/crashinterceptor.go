package serverinterceptors

import (
	"context"
	"runtime/debug"

	"github.com/WeiXinao/daily_your_go/pkg/log"
	"google.golang.org/grpc"
)

func StreamCrashInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer handleCrash(func(a any) {
		log.Errorf("%+v\n\t%s", a, debug.Stack())
	})

	return handler(srv, ss)
}

func UnaryCrashInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	defer handleCrash(func(a any) {
		log.Errorf("%+v\n\t%s", a, debug.Stack())
	})

	return handler(ctx, req)
}

func handleCrash(handler func(any)) {
	if err := recover(); err != nil {
		handler(err)
	}
}