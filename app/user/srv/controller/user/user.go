package user

import srvv1 "github.com/WeiXinao/daily_your_go/app/user/srv/service/v1"

type userServer struct {
	srv srvv1.UserSrv
}

// java 中的 ioc，@Autowire 控制反转 ioc 
// 代码分层，第三方服务，rpc，redis 等等，带来了一定的复杂度
func NewUserServer(srv srvv1.UserSrv) *userServer {
	return &userServer{srv: srv}
}
