package infra

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewNacosConfigClient,
	InitNacosSubscriber,
	NewRegistrar,
	NewGormDB,
)