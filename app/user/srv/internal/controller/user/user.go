package user

import (

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	upbv1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	srvv1 "github.com/WeiXinao/daily_your_go/app/user/srv/internal/service/v1"
)

var _ upb.UserServer = (*userServer)(nil)

type userServer struct {
	srv srvv1.UserSrv
	upb.UnimplementedUserServer
}

// java 中的 ioc，@Autowire 控制反转 ioc
// 代码分层，第三方服务，rpc，redis 等等，带来了一定的复杂度
func NewUserServer(srv srvv1.UserSrv) upb.UserServer {
	return &userServer{srv: srv}
}

func DTOToResponse(userdto srvv1.UserDTO) *upbv1.UserInfoResponse {
	// 在 grpc 的 message 中字段有默认值，你不能随便赋值 nil 进去，容易出错
	// 这里要搞清楚，哪些字段有默认值
	userInfoRsp := &upbv1.UserInfoResponse{
		Id:       userdto.ID,
		Password: userdto.Password,
		NickName: userdto.NickName,
		Mobile:   userdto.Mobile,
		Gender:   userdto.Gender,
		Role:     int32(userdto.Role),
	}

	if userdto.Birthday != nil {
		userInfoRsp.Birthday = uint64(userdto.Birthday.Unix())
	}
	return userInfoRsp	
}