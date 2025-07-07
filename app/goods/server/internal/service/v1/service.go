package service

import (
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	search "github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1"
)

type ServiceFactory interface {
	Goods() GoodsSvc
}

var _ ServiceFactory = (*service)(nil)

type service struct {
	data       data.DataFactory
	dataSearch search.SearchFactory
}

func NewServiceFactory(dataFactory data.DataFactory, searchFactory search.SearchFactory) ServiceFactory {
	return &service{
		data:       dataFactory,
		dataSearch: searchFactory,
	}
}

// Goods implements Service.
func (s *service) Goods() GoodsSvc {
	return newGoodsService(s.data, s.dataSearch)
}
