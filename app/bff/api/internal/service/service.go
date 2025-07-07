package service

import (
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/sms/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service/user/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
)

type ServiceFactory interface {
	Goods() goods.GoodsSrv
	Users() user.UserSrv
	Sms() sms.SmsSrv
}

var _ ServiceFactory = (*service)(nil)

type service struct {
	data data.DataFactory

	smsOpts *options.SmsOptions
	jwtOpts *options.JwtOptions
}

// Sms implements ServiceFactory.
func (s *service) Sms() sms.SmsSrv {
	return sms.NewSmsService(s.smsOpts)
}

// Users implements ServiceFactory.
func (s *service) Users() user.UserSrv {
	return user.NewUserService(s.data, s.jwtOpts)
}

// Goods implements Service.
func (s *service) Goods() goods.GoodsSrv {
	return goods.NewGoods(s.data)
}

func NewServiceFactory(
		dataFactory data.DataFactory, 
		smsOpts *options.SmsOptions, 
		jwtOpts *options.JwtOptions,
	) ServiceFactory {
	return &service{
		data:    dataFactory,
		smsOpts: smsOpts,
		jwtOpts: jwtOpts,
	}
}
