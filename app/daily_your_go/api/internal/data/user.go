package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/pkg/common/time"
)

type User struct {
	ID       uint64     `json:"id,omitempty"`
	Mobile   string     `json:"mobile,omitempty"`
	Password string     `json:"password,omitempty"`
	NickName string     `json:"nick_name,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	Gender   string     `json:"gender,omitempty"`
	Role     int        `json:"role,omitempty"`
}

type UserList struct {
	TotalCount int64   `json:"totalCount"`
	Items      []*User `json:"items"`
}

type UserData interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Get(ctx context.Context, id uint64) (*User, error)
	GetByMobile(ctx context.Context, mobile string) (*User, error)
}