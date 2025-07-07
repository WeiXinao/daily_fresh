package user

import (
	"context"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

// GetUserById implements v1.UserServer.
func (u *userServer) GetUserById(ctx context.Context, request *upb.IdRequest) (*upb.UserInfoResponse, error) {
	log.Infof("get user by id function called.")
	user, err := u.srv.GetByID(ctx, uint64(request.GetId()))
	if err != nil {
		log.Errorf("get user by id: %s, error: %v", request.GetId(), err)
		return nil, err
	}
	
	return DTOToResponse(*user), nil
}