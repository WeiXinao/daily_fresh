package mock

import (
	"context"

	metav1 "github.com/WeiXinao/daily_fresh/pkg/common/meta/v1"
	udv1 "github.com/WeiXinao/daily_fresh/app/user/server/internal/data/v1"
)

type users struct {
	users []*udv1.UserDO
}

func NewUsers() *users {
	return &users{}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*udv1.UserDOList, error) {
	usrs := []*udv1.UserDO{
		{NickName: "xiaoxin"},
	}
	return &udv1.UserDOList{
		TotalCount: 1,
		Items: usrs,	
	}, nil
}

