package v1

import (
	"context"
	"testing"

	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/app/user/srv/data/v1/mock"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), metav1.ListMeta{})
}