package app

import (
	"context"
	"maps"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts options

	lk sync.RWMutex

	instance *registry.ServiceInstance
	bgCtxWithCancel context.Context
	cancel func()
	shutdownChan chan struct{}
	finishChan chan struct{}
}

func New(opts ...Option) *App {
	o := options{
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registerTimeout: 10 * time.Second,
		stopTimeout:     10 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
		finishChan: make(chan struct{}, 1),
	}
}

// 启动整个微服务
func (a *App) Run() error {
	// 注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	// 这个变量可能被其他的 goroutine 访问到
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	// 重点，写的很简单，http 服务要启动
	// 现在启动了两个 server，一个是 restserver，一个是 rpcserver
	/*
		1. 这两个 server 是否必须同时启动成功
		2. 如果有一个启动失败，那么我们应该停止另外一个 server
		3. 如果启动了多个，如果其中一个启动失败，其他的应该被取消
			如果剩余的server的状态
			1. 还没有开始调用 start  不进行就行了，stop 也行
			2. start 进行中  调用进行中的 cancel
			3. start 已经完成 调用 stop
		如果我们的服务启动了然后这个时候用户立马进行了访问
	*/

	servers := []server.Server{}
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}

	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}

	bgCtxWithCancel, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	eg, ctx := errgroup.WithContext(bgCtxWithCancel)
	wg := sync.WaitGroup{}
	for _, srv := range servers {
		// 启动 server
		// 在启动一个 goroutine 去监听是否有 err 产生
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal

			// 不可能无休止的等待 stop
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})

		wg.Add(1)
		eg.Go(func () error {
			wg.Done()
			log.Info("start server")
			return srv.Start(ctx)
		})
	}

	wg.Wait()

	// 注册服务
	if a.opts.registrar != nil {
		rctx, rconcel := context.WithTimeout(context.Background(), a.opts.registerTimeout)
		defer rconcel()
		err := a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error: %s", err)
			return err
		}
	}

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	a.shutdownChan = make(chan struct{})
	signal.Notify(quit, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				return a.Stop()
			case <-a.shutdownChan:
				return a.Stop()
			}
		}
	})
 	err = eg.Wait(); 
	if err != nil {
		return err
	}
	a.finishChan<-struct{}{}

	return nil
}

/*
http basic 认证
cache：1. redis 2. memcache 3. local cache
*/
// 停止服务
func (a *App) Stop() error {
	a.lk.RLock()
	insecure := a.instance
	a.lk.RUnlock()

	log.Info("start deregister service")
	if a.instance != nil && a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, insecure); err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}

	if a.cancel != nil {
		log.Info("start cancel context")
		a.cancel()
	}

	return nil
}

func (a *App) Shutdown() error {
	close(a.shutdownChan)
	<-a.finishChan
	return nil
}

// 创建服务注册的结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	metadata := make(map[string]string, len(a.opts.metadata))
	maps.Copy(metadata, a.opts.metadata)

	// 默认使用 id 作为 hashKey
	metadata["hashKey"] = a.opts.id

	// 从 rpcserver，restserver 去主动获取这些信息
	if a.opts.rpcServer != nil {
		u := &url.URL{
			Scheme: "grpc",
			Path: a.opts.rpcServer.Address(),
		}
		endpoints = append(endpoints, u.String())
	}

	if a.opts.restServer != nil {
		u := &url.URL{
			Scheme: "http",
			Path: a.opts.restServer.Address(),
		}
		endpoints = append(endpoints, u.String())
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
		Metadata: metadata,
	}, nil
}
