package user

import (
	"context"
	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	srvv1 "github.com/WeiXinao/daily_your_go/app/user/srv/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ upb.UserServer = (*userServer)(nil)

type userServer struct {
	srv srvv1.UserSrv
	upb.UnimplementedUserServer
}

// CheckPassword implements v1.UserServer.
func (u *userServer) CheckPassword(context.Context, *upb.PasswordCheckInfo) (*upb.CheckResponse, error) {
	panic("unimplemented")
}

// CreateUser implements v1.UserServer.
func (u *userServer) CreateUser(context.Context, *upb.CreateUserInfo) (*upb.UserInfoResponse, error) {
	panic("unimplemented")
}

// GetUserById implements v1.UserServer.
func (u *userServer) GetUserById(context.Context, *upb.IdRequest) (*upb.UserInfoResponse, error) {
	panic("unimplemented")
}

// GetUserByMobile implements v1.UserServer.
func (u *userServer) GetUserByMobile(context.Context, *upb.MobileRequest) (*upb.UserInfoResponse, error) {
	panic("unimplemented")
}

// UpdateUser implements v1.UserServer.
func (u *userServer) UpdateUser(context.Context, *upb.UpdateUserInfo) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// java 中的 ioc，@Autowire 控制反转 ioc
// 代码分层，第三方服务，rpc，redis 等等，带来了一定的复杂度
func NewUserServer(srv srvv1.UserSrv) *userServer {
	return &userServer{srv: srv}
}
