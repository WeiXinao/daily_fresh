package srv

import (
	"fmt"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/user/srv/config"
	"github.com/WeiXinao/daily_your_go/app/user/srv/controller/user"
	"github.com/WeiXinao/daily_your_go/app/user/srv/data/v1/mock"
	svcv1 "github.com/WeiXinao/daily_your_go/app/user/srv/service/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
)

func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 有点繁琐 wire, ioc-golang
	data := mock.NewUsers()
	svc := svcv1.NewUserService(data)
	usrv := user.NewUserServer(svc)
	

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))	
	upb.RegisterUserServer(urpcServer.Server, usrv)
	return nil, nil
}