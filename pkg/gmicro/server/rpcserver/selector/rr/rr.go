package rr

import (
	"context"
	"sync/atomic"

	selector2 "github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector/node/direct"
)

const (
	Name = "rr"
)

func New() selector2.Selector {
	return NewBuilder().Build()
}

func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Balancer: &builder{},
		Node: &direct.Builder{},
	}
}

var _ selector2.BalancerBuilder = (*builder)(nil)

type builder struct{}

// Build implements selector.BalancerBuilder.
func (b *builder) Build() selector2.Balancer {
	balancer := &balancer{
		idx: atomic.Value{},
	}
	balancer.idx.Store(0)
	return balancer
}

var _ selector2.Balancer = (*balancer)(nil)

type balancer struct {
	idx atomic.Value
}

// Pick implements selector.Balancer.
func (b *balancer) Pick(ctx context.Context, nodes []selector2.WeightedNode) (selected selector2.WeightedNode, done selector2.DoneFunc, err error) {
	// 重试 3 次，使负载尽量均衡
	for range 3 {
		var i = b.idx.Load().(int) 
		selected = nodes[i]
		swapped := b.idx.CompareAndSwap(i, i+1)
		if swapped {
			break
		}
	}
	d := selected.Pick()
	return selected, d, nil	
}
