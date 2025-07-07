package srv

import (
	"fmt"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/config"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/controller/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1/db"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1/es"
	svcv1 "github.com/WeiXinao/daily_your_go/app/goods/srv/internal/service/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/core/trace"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

func NewGoodsRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 opentelemetry 的 exporter
	trace.InitAgent(trace.Options{
		Name: cfg.Telemtry.Name,
		Endpoint: cfg.Telemtry.Endpoint,
		Batcher: cfg.Telemtry.Batcher,
		Sampler: cfg.Telemtry.Sampler,
	})

	// 构建，繁琐 - 工厂模式
	// 有点繁琐 wire, ioc-golang
	dataFactory, err := db.GetDBFactoryOr(cfg.MySQL)
	if err!= nil {
		log.Fatal(err.Error())
	}

	searchFactory, err := es.GetSearchFactoryOr(cfg.Es)
	if err != nil {
		log.Fatal(err.Error())
	}

	// ioc 框架 wire，ioc-golang（处于早期，有坑）
	// 基于工厂方法
	svcFactory := svcv1.NewServiceFactory(dataFactory ,searchFactory)
	gctrl := controller.NewGoodServer(svcFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))	

	gpb.RegisterGoodsServer(grpcServer.Server, gctrl)
	
	return grpcServer, nil
}