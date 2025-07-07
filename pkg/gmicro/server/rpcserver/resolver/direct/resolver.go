package direct

import "google.golang.org/grpc/resolver"

var _ resolver.Resolver = (*directResolver)(nil)

func init() {
	resolver.Register(NewBuilder())
}

type directResolver struct {
}

func newDirectResolver() *directResolver {
	return &directResolver{}
}

// Close implements resolver.Resolver.
func (d *directResolver) Close() {
}

// ResolveNow implements resolver.Resolver.
func (d *directResolver) ResolveNow(resolver.ResolveNowOptions) {
}
