package owrr

import (
	"context"
	"sync"

	selector2 "github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector/node/direct"
	"github.com/WeiXinao/daily_fresh/pkg/errors"
)

const (
	Name = "owrr"
)

var (
	ErrInitWeight = errors.New("someone's weight is not inited")
)

func New() selector2.Selector {
	return NewBuilder().Build()
}

var _ selector2.Balancer = (*balancer)(nil)
type balancer struct {
	mu sync.Mutex
	total int64
	idx int
	curTotalWeight int64
}

// Pick implements selector.Balancer.
func (b *balancer) Pick(ctx context.Context, nodes []selector2.WeightedNode) (selected selector2.WeightedNode, done selector2.DoneFunc, err error) {
	var (
		g = *nodes[0].InitialWeight()
		totalWeight = int64(0)
	)
	for _, node := range nodes {
		if node.InitialWeight() == nil {
			return nil, nil, ErrInitWeight
		}
		iw := *node.InitialWeight()
		totalWeight += iw
		g = gcd(iw, g)
	}

	b.mu.Lock()
	b.total = (b.total + g) % totalWeight
	if b.total > b.curTotalWeight {
		b.idx = (b.idx + 1) % len(nodes)
		b.curTotalWeight = (b.curTotalWeight + *nodes[b.idx].InitialWeight()) % totalWeight
	}
	b.mu.Unlock()

	selected = nodes[b.idx]
	d := selected.Pick()
	return selected, d, nil
}

func gcd(a, b int64) int64 {
	if a < b {
		a, b = b, a
	}
	for b != 0 {
		a, b = b, a % b
	}
	return a
}

var _ selector2.BalancerBuilder = (*builder)(nil)

type builder struct{}

// Build implements selector.BalancerBuilder.
func (b *builder) Build() selector2.Balancer {
	return &balancer{}
}

func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Balancer: &builder{},
		Node: &direct.Builder{},
	}
}
