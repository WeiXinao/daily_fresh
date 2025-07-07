package server

import (
	"context"
	"github.com/WeiXinao/daily_fresh/test/gin_grpc_test/proto"
)

var _ proto.HelloServiceServer = (*sayHelloServer)(nil)

type sayHelloServer struct {
	proto.UnimplementedHelloServiceServer
}

func NewSayHelloServer() proto.HelloServiceServer {
	return &sayHelloServer{}
}

// SayBye implements proto.HelloServiceServer.
func (s *sayHelloServer) SayBye(context.Context, *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	panic("unimplemented")
}

// SayHello implements proto.HelloServiceServer.
func (s *sayHelloServer) SayHello(ctx context.Context, request *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	return &proto.SayHelloResponse{
		Greet: "hello! " + request.GetName(),
	}, nil
}
