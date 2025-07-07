package service

import (
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/data/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
)

type ServiceFactory interface {
	Order() OrderSrv
}

var _ ServiceFactory = (*service)(nil)

type service struct {
	data    data.DataFactory
	dtmOpts *options.DtmOptions
}

// Order implements ServiceFactory.
func (s *service) Order() OrderSrv {
	return newOrderService(s.data, s.dtmOpts)
}

func NewServiceFactory(data data.DataFactory, dtmOpts *options.DtmOptions) ServiceFactory {
	return &service{
		data:    data,
		dtmOpts: dtmOpts,
	}
}
