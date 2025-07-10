package user

import (
	"context"

	upb "github.com/WeiXinao/daily_fresh/api/user/v1"
	v1 "github.com/WeiXinao/daily_fresh/app/user/server/internal/service/v1"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/code"
	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/jinzhu/copier"
)

// CreateUser implements v1.UserServer.
func (u *userServer) CreateUser(ctx context.Context,request *upb.CreateUserInfo) (*upb.UserInfoResponse, error) {
	log.Infof("create user function is called.")
	
	userDTO := v1.UserDTO{}	
	err := copier.Copy(&userDTO, request)
	if err != nil {
		log.Errorf("request 转化为 userDTO 失败：%v", err.Error())
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}
	
	err = u.srv.Create(ctx, &userDTO)
	if err != nil {
		log.Errorf("create user: %v, error: %v", userDTO, err.Error())
		return nil, err
	}

	return DTOToResponse(userDTO), nil
}