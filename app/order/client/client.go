package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	proto "github.com/WeiXinao/daily_your_go/api/order/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/registry/consul"
	rpc "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	_ "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/resolver/direct"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/selector"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/selector/random"
	"github.com/hashicorp/consul/api"
)

func generateOrderSn(userId int32) string {
	/*
	 订单号的生成规则：
	 年月日时分秒 + 用户id + 2 位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

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
		rpc.WithEndpoint("discovery:///daily-your-go-order-srv"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userId := 1
	oc := proto.NewOrderClient(conn)
	_, err = oc.SubmitOrder(context.Background(), &proto.OrderRequest{
		UserId: int32(userId),
		Address: "慕课网",
		OrderSn: generateOrderSn(int32(userId)),
		Name: "小新",
		Post: "请尽快发货",
		Mobile: "18787878787",
	})
	if err != nil {
		panic(err)
	}
}