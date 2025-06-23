package srv

import (
	"fmt"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/core/trace"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
)

func NewUserRPCServer(
	telemetryOpts *options.TelemtryOptions,
	serverOpts *options.ServerOptions,
	usrv upb.UserServer,
) (*rpcserver.Server, error) {
	trace.InitAgent(trace.Options{
		Name: telemetryOpts.Name,
		Endpoint: telemetryOpts.Endpoint,
		Batcher: telemetryOpts.Batcher,
		Sampler: telemetryOpts.Sampler,
	})

	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))	

	upb.RegisterUserServer(urpcServer.Server, usrv)
	
	return urpcServer, nil
}
