package service

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/user/v1"
)

type ServiceFactory interface {
	Goods() goods.GoodsSrv
	Users() user.UserSrv
}