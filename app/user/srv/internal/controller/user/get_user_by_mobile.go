package user

import (
	"context"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/daily_your_go/pkg/private"
)

// GetUserByMobile implements v1.UserServer.
func (u *userServer) GetUserByMobile(ctx context.Context, request *upb.MobileRequest) (*upb.UserInfoResponse, error) {
	log.Infof("get user by mobile function called.")
	user, err := u.srv.GetByMobile(ctx, request.GetMobile())
	if err != nil {
		log.Errorf("get user by mobile: %s, error: %v", private.InsensitiveMobile(request.GetMobile()), err)
		return nil, err
	}
	
	return DTOToResponse(*user), nil
}