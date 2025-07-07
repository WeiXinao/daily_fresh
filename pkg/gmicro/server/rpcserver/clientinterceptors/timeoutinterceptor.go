package clientinterceptors

import (
	"context"
	"sync"
	"time"

	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func (ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx, cancal := context.WithTimeout(ctx, timeout)	
		defer cancal()
		var (
			panicChan = make(chan any, 1)
			doneChan = make(chan struct{})
			err error
			mx sync.Mutex
		)
		go func ()  {
			defer func ()  {
				if e := recover(); e != nil {
					panicChan <- errors.Errorf("%+v", e)
				}
			}()
			mx.Lock()
			err = invoker(ctx, method, req, reply, cc, opts...)
			mx.Unlock()
			close(doneChan)
		}()
		
		select {
		case p := <- panicChan:
			panic(p)
		case <-doneChan:
			mx.Lock()
			defer mx.Unlock()	
			return err
		case <-ctx.Done():
			er := ctx.Err()
			if errors.Is(ctx.Err(), context.Canceled) {
				er = status.Error(codes.Canceled, er.Error())
			} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				er = status.Error(codes.DeadlineExceeded, er.Error())
			}
			return er
		}
	}
}