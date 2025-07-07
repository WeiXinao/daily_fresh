package srv

import (
	"fmt"

	proto "github.com/WeiXinao/daily_your_go/api/order/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/config"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/controller/order/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/data/v1/db"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/service/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/core/trace"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

func NewOrderRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 opentelemetry 的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemtry.Name,
		Endpoint: cfg.Telemtry.Endpoint,
		Batcher:  cfg.Telemtry.Batcher,
		Sampler:  cfg.Telemtry.Sampler,
	})

	// 构建，繁琐 - 工厂模式
	// 有点繁琐 wire, ioc-golang
	dataFactory, err := db.GetDataFactory(cfg.MySQL, cfg.Registry)
	if err != nil {
		log.Fatal(err.Error())
	}

	// ioc 框架 wire，ioc-golang（处于早期，有坑）
	// 基于工厂方法
	svcFactory := service.NewServiceFactory(dataFactory, cfg.Dtm)
	octrl := order.NewOrderServer(svcFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	proto.RegisterOrderServer(grpcServer, octrl)

	return grpcServer, nil
}
