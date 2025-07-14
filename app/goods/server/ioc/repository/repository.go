package repository

import (
	"github.com/WeiXinao/daily_fresh/app/goods/server/ioc/repository/dao"
	"github.com/WeiXinao/daily_fresh/app/goods/server/ioc/repository/es"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(dao.ProviderSet, es.ProviderSet)