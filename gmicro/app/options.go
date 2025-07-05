package app

import (
	"net/url"
	"os"
	"time"

	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
)

type Option func(o *options)

type options struct {
	id        string
	name      string
	endpoints []*url.URL
	metadata map[string]string

	sigs []os.Signal

	// 允许用户传入自己的实现
	registrar       registry.Registrar
	registerTimeout time.Duration
	
	// stop 超时时间
	stopTimeout     time.Duration

	restServer *restserver.Server
	rpcServer *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}

func WithRestServer(server *restserver.Server) Option {
	return func(o *options) {
		o.restServer = server
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithMetadata(key, value string) Option {
	return func(o *options) {
		o.metadata[key] = value
	}
}

func WithSignals(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
