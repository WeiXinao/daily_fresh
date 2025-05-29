package main

import (
	"context"
	"fmt"

	v1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	rpc "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	_ "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/resolver/direct"
	"github.com/hashicorp/consul/api"
)

func main() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.5.52:8500"
	cfg.Scheme = "http"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	conn, err := rpc.DialInsecure(
		context.Background(), 
		rpc.WithDiscovery(r),
		rpc.WithEndpoint("discovery:///daily-your-go-user-srv"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)
	ulr, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
	if err != nil {
		panic(err)
	}
	fmt.Println(ulr)
}