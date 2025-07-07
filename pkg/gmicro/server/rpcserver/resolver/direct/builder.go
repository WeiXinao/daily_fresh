package direct

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*directBuilder)(nil)

type directBuilder struct{}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
// example:
//	
// direct://<authority>/127.0.0.1:9000

func NewBuilder() *directBuilder {
	return &directBuilder{}
}

// Build implements resolver.Builder.
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{ Addr: addr })
	}

	// grpc 建立连接的逻辑都在这里
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

// Scheme implements resolver.Builder.
func (d *directBuilder) Scheme() string {
	return "direct"
}
