package srv

import (
	"fmt"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1/es"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/config"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1/db"
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

	// 有点繁琐 wire, ioc-golang
	gdb, err := db.GetDBFactoryOr(cfg.MySQL)
	if err!= nil {
		log.Fatal(err.Error())
	}

	c, err := es.GetSearchFactoryOr(cfg.Es)
	if err != nil {
		log.Fatal(err.Error())
	}

	data := db.NewGoods(gdb)
	svc := svcv1.NewUserService(data)
	usrv := user.NewUserServer(svc)
	

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))	

	upb.RegisterUserServer(urpcServer.Server, usrv)
	
	return urpcServer, nil
}