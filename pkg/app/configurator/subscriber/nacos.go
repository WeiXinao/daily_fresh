package subscriber

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var _ Subscriber = (*NacosSubscriber)(nil)

type NacosSubscriber struct {
	nacosConfigClient config_client.IConfigClient
	dataId string
	group string
}

func NewNacosSubscriber(nacosConfigClient config_client.IConfigClient, 
	group, dataId string) *NacosSubscriber {
	return &NacosSubscriber {
		nacosConfigClient: nacosConfigClient,
		group: group,
		dataId: dataId,
	}
}

// AddListener implements Subscriber.
func (n *NacosSubscriber) AddListener(listener func(key string, data string)) error {
	return n.nacosConfigClient.ListenConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group: n.group,
		OnChange: func(namespace string, group string, dataId string, data string) {
			listener(dataId, data)
		},
	})
}

// Value implements Subscriber.
func (n *NacosSubscriber) Value() (string, error) {
	return n.nacosConfigClient.GetConfig(vo.ConfigParam{
    DataId: n.dataId,
    Group:  n.group,
	})
}

func (n *NacosSubscriber) Key() string {
	return n.dataId
}

