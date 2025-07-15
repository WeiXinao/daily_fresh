package infra

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func NewNacosConfigClient(nacos *options.NacosOptions) (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: nacos.Host,
			Port: uint64(nacos.Port),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         nacos.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, err
	}
	return configClient, nil
}

func InitNacosSubscriber(configClient config_client.IConfigClient, nacosOpts *options.NacosOptions) (*subscriber.NacosSubscriber, error) {
	return subscriber.NewNacosSubscriber(configClient, nacosOpts.Group, nacosOpts.DataId), nil
}