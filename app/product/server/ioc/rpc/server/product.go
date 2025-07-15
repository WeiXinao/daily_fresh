package server

import (
	"fmt"

	gpb "github.com/WeiXinao/daily_fresh/api/goods/v1"
	proto "github.com/WeiXinao/daily_fresh/api/goods/v1"
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/core/trace"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/server/rpcserver"
)

func NewProductRPCServer(
	telemetryOpts *options.TelemtryOptions,
	serverOpts *options.ServerOptions,
	gsrv proto.GoodsServer,
) (*rpcserver.Server, error) {
	// 初始化 opentelemetry 的 exporter
	trace.InitAgent(trace.Options{
		Name: telemetryOpts.Name,
		Endpoint: telemetryOpts.Endpoint,
		Batcher: telemetryOpts.Batcher,
		Sampler: telemetryOpts.Sampler,
	})

	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))	

	gpb.RegisterGoodsServer(grpcServer.Server, gsrv)
	
	return grpcServer, nil
}