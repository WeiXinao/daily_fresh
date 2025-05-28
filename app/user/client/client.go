package main

import (
	"context"
	"fmt"

	v1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"127.0.0.1:8078", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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