package rpc

import (
	"github.com/WeiXinao/daily_fresh/app/goods/server/ioc/rpc/client"
	"github.com/WeiXinao/daily_fresh/app/goods/server/ioc/rpc/server"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(client.ProviderSet, server.ProviderSet)