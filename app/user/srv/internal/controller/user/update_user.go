package user

import (
	"context"
	"time"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	v1 "github.com/WeiXinao/daily_your_go/app/user/srv/internal/service/v1"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser implements v1.UserServer.
func (u *userServer) UpdateUser(ctx context.Context, request *upb.UpdateUserInfo) (*emptypb.Empty, error) {
	log.Info("update user function called.")

	userDTO := &v1.UserDTO{
	}
	err := copier.Copy(&userDTO, request)
	if err != nil {
		log.Errorf("request 转化为 userDTO 失败：%v", err.Error())
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}
	birthday := time.Unix(int64(request.Birthday), 0)
	userDTO.Birthday =  &birthday

	err = u.srv.Update(ctx, userDTO)	
	if err != nil {
		log.Errorf("update user: %v, error: %v", userDTO, err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}