package rpc

import (
	"context"
	"time"

	upbv1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/jinzhu/copier"
	timePkg "github.com/WeiXinao/daily_your_go/pkg/common/time"
)

var _ data.UserData = (*users)(nil)

type users struct {
	uc upbv1.UserClient
}

// Create implements data.UserData.
func (u *users) Create(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.CreateUserInfo{
		NickName: user.NickName,
		Mobile: user.Mobile,
		Password: user.Password,
	}

	uir, err := u.uc.CreateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	user.ID = uint64(uir.Id)
	return nil
}

// Get implements data.UserData.
func (u *users) Get(ctx context.Context, id uint64) (*data.User, error) {
	user, err := u.uc.GetUserById(ctx, &upbv1.IdRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}
	
	userDO := &data.User{}
	err = copier.Copy(userDO, user)
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, "proto message转化为DO错误，err:%w", err)
	}
	userDO.Birthday = timePkg.Time{Time: time.Unix(int64(user.Birthday), 0)}
	return userDO, nil
}

// GetByMobile implements data.UserData.
func (u *users) GetByMobile(ctx context.Context, mobile string) (*data.User, error) {
	user, err := u.uc.GetUserByMobile(ctx, &upbv1.MobileRequest{Mobile: mobile})
	if err != nil {
		return nil, err
	}
	
	userDO := &data.User{}
	err = copier.Copy(userDO, user)
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, "proto message转化为DO错误，err:%w", err)
	}
	userDO.Birthday = timePkg.Time{Time: time.Unix(int64(user.Birthday), 0)}
	return userDO, nil
}

// Update implements data.UserData.
func (u *users) Update(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.UpdateUserInfo{}
	err := copier.Copy(protoUser, user)
	if err != nil {
		return errors.WithCode(code.ErrUnknown, "DO转化为proto message错误，err:%w", err)
	}
	protoUser.Birthday = uint64(user.Birthday.Unix())

	_, err = u.uc.UpdateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	return nil
}
