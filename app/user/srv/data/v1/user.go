package v1

import (
	"context"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
)

type UserStore interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDOList, error)
}

type UserDO struct {
	Name string
}

type UserDOList struct {
	TotalCount int64 // 总数
	Items []*UserDO 
}
