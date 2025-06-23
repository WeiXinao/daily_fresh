package srv

import (
	"fmt"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/core/trace"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func NewNacosDatasource(nacosOptions *options.NacosOptions) (*nacos.NacosDataSource, error) {
		// nacos server 地址
	sc := []constant.ServerConfig{
		{
			ContextPath: "/nacos",
			IpAddr:      nacosOptions.Host,
			Port:        uint64(nacosOptions.Port),
		},
	}

	// nacos client 相关的参数配置，，具体配置参考 github.com/WeiXinao/StudyCode/sentinel_nacos
	cc := constant.ClientConfig{
		NamespaceId: nacosOptions.Namespace,
		TimeoutMs: 5000,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig": cc,
	})
	if err != nil {
		return nil, err
	}

	// 注册流控规则 Handler
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)

	// 创建 NacosDataSource 的数据源
	nds, err := nacos.NewNacosDataSource(client, nacosOptions.Group, nacosOptions.DataId, h)
	if err != nil {
		return nil, err
	}
	return nds, nil
}

func NewUserRPCServer(
	telemetryOpts *options.TelemtryOptions,
	serverOpts *options.ServerOptions,
	usrv upb.UserServer,
	nds *nacos.NacosDataSource,
) (*rpcserver.Server, error) {
	// 初始化 open-telemetry 的 exporter
	trace.InitAgent(trace.Options{
		Name: telemetryOpts.Name,
		Endpoint: telemetryOpts.Endpoint,
		Batcher: telemetryOpts.Batcher,
		Sampler: telemetryOpts.Sampler,
	})


	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)
	opts := []rpcserver.ServerOption{
		rpcserver.WithAddress(rpcAddr),
	}
	if serverOpts.EnableLimit {
		err := nds.Initialize()
		if err != nil {
			return nil, err
		}
		opts = append(opts, rpcserver.WithUnaryInterceptor(grpc.NewUnaryServerInterceptor()))
	}
	urpcServer := rpcserver.NewServer(opts...)	

	upb.RegisterUserServer(urpcServer.Server, usrv)
	
	return urpcServer, nil
}
