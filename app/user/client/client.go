package main

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/WeiXinao/daily_fresh/api/user/v1"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/registry/consul"
	rpc "github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver"
	_ "github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/resolver/direct"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver/selector/random"
	"github.com/hashicorp/consul/api"
)

func main() {
	// 设置全局的负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()

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
		rpc.WithBalacerName("selector"),
		rpc.WithDiscovery(r),
		rpc.WithClientTimeout(1 * time.Hour),
		rpc.WithEndpoint("discovery:///daily-your-go-user-srv"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)

	for {
		_, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
		if err != nil {
			panic(err)
		}
		fmt.Println("success")
		time.Sleep(time.Millisecond * 2)
	}
}