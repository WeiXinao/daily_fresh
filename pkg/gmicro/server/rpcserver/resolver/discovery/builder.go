package discovery

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry"
	"google.golang.org/grpc/resolver"
)

const name = "discovery"

// Option is builder option.
type Option func(o *builder)


// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
    return func(b *builder) {
        b.timeout = timeout
    }
}

// WithInsecure with isSecure option.
func WithInsecure(insecure bool) Option {
    return func(b *builder) {
        b.insecure = insecure
    }
}

type builder struct {
    discoverer registry.Discovery
    timeout    time.Duration
    insecure   bool
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
    b := &builder{
        discoverer: d,
        timeout:    time.Second * 10,
        insecure:   false,
    }
    for _, o := range opts {
        o(b)
    }
    return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    var (
        err error
        w   registry.Watcher
    )
    done := make(chan struct{}, 1)
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        w, err = b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
        close(done)
    }()
    select {
    case <-done:
    case <-time.After(b.timeout):
        err = errors.New("discovery create watcher overtime")
    }
    if err != nil {
        cancel()
        return nil, err
    }
    r := &discoveryResolver{
        w:        w,
        cc:       cc,
        ctx:      ctx,
        cancel:   cancel,
        insecure: b.insecure,
    }
    go r.watch()
    return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
    return name
}