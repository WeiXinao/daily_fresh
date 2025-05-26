package db

import (
	"context"

	"gorm.io/gorm"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	udv1 "github.com/WeiXinao/daily_your_go/app/user/srv/data/v1"
)

var _ udv1.UserStore = (*users)(nil)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*udv1.UserDOList, error) {
	// 你要用 gorm 还是其他的 orm 无所谓
	return &udv1.UserDOList{}, nil
}