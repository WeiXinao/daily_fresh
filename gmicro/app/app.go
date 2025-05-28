package app

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/google/uuid"
)

type App struct {
	opts options

	lk sync.RWMutex

	instance *registry.ServiceInstance
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

	return &App{opts: o}
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
	if a.opts.rpcServer != nil {
		go func() {
			err := a.opts.rpcServer.Start()
			if err != nil {
				panic(err)
			}
		}()
	}

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
	signal.Notify(quit, a.opts.sigs...)
	<-quit

	return nil
}

// 停止服务
func (a *App) Stop() error {
	a.lk.RLock()
	insecure := a.instance
	a.lk.RUnlock()

	if a.instance != nil && a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, insecure); err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}
	return nil
}

// 创建服务注册的结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	// 从 rpcserver，restserver 去主动获取这些信息
	if a.opts.rpcServer != nil {
		u := &url.URL{
			Scheme: "grpc",
			Path: a.opts.rpcServer.Address(),
		}
		endpoints = append(endpoints, u.String())
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
