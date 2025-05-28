package main

import (
	"context"
	"fmt"

	v1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	rpc "github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
)

func main() {
	conn, err := rpc.DialInsecure(
		context.Background(), 
		rpc.WithEndpoint("192.168.5.52:8078"),
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