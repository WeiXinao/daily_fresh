package user

import (
	"context"

	upbv1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	srvv1 "github.com/WeiXinao/daily_your_go/app/user/srv/service/v1"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
)

func DTOToResponse(userdto srvv1.UserDTO) *upbv1.UserInfoResponse {
	return &upbv1.UserInfoResponse{NickName: userdto.Name}
}

/*
controller 层依赖了 service，service 层依赖了 data 层：
	controller 层能否直接依赖 data 层：可以的
controller 依赖 service 并不是直接依赖了具体的 struct，而是依赖了 interface
*/
func (u *userServer) GetUserList(ctx context.Context, req *upbv1.PageInfo) (*upbv1.UserListResponse, error) {
	srvOpts := metav1.ListMeta{
		Page: int(req.Pn),
		PageSize: int(req.PSize),
	}
	dtoList, err := u.srv.List(ctx, srvOpts) 
	if err != nil {
		return nil, err
	}

	var rsp upbv1.UserListResponse
	for _, value := range dtoList.Items {
		rsp.Data = append(rsp.Data, DTOToResponse(*value)) 	
	}
	return &rsp, nil
}