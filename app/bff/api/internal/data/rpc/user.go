package rpc

import (
	"context"
	"time"

	upbv1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	v1 "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/gmicro/registry"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver"
	"github.com/WeiXinao/daily_your_go/gmicro/server/rpcserver/clientinterceptors"
	timePkg "github.com/WeiXinao/daily_your_go/pkg/common/time"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/jinzhu/copier"
)

var _ data.UserData = (*users)(nil)

type users struct {
	uc upbv1.UserClient
}

// CheckPassword implements data.UserData.
func (u *users) CheckPassword(ctx context.Context, pwd string, encryptedPwd string) (bool, error) {
	rsp, err := u.uc.CheckPassword(ctx, &upbv1.PasswordCheckInfo{
		Password: pwd,
		EncryptedPassword: encryptedPwd,
	})
	if err != nil {
		return false, err
	}
	return rsp.GetSuccess(), nil
}

// Create implements data.UserData.
func (u *users) Create(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.CreateUserInfo{
		NickName: user.NickName,
		Mobile:   user.Mobile,
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
func (u *users) Get(ctx context.Context, id uint64) (data.User, error) {
	user, err := u.uc.GetUserById(ctx, &upbv1.IdRequest{Id: int32(id)})
	if err != nil {
		return data.User{}, err
	}

	userDO := data.User{}
	err = copier.Copy(&userDO, user)
	if err != nil {
		return data.User{}, errors.WithCode(code.ErrUnknown, "proto message转化为DO错误，err:%w", err)
	}
	userDO.Birthday = timePkg.Time{Time: time.Unix(int64(user.Birthday), 0)}
	return userDO, nil
}

// GetByMobile implements data.UserData.
func (u *users) GetByMobile(ctx context.Context, mobile string) (data.User, error) {
	user, err := u.uc.GetUserByMobile(ctx, &upbv1.MobileRequest{Mobile: mobile})
	if err != nil {
		return data.User{}, err
	}

	userDO := data.User{}
	err = copier.Copy(&userDO, user)
	if err != nil {
		return data.User{}, errors.WithCode(code.ErrUnknown, "proto message转化为DO错误，err:%w", err)
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

const userServiceName = "discovery:///daily-your-go-user-srv"

func NewUserServiceClient(r registry.Discovery) upbv1.UserClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(), 
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientTimeout(1 * time.Hour),
		rpcserver.WithEndpoint(userServiceName),
		rpcserver.WithClientUnaryInterceptors(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	return v1.NewUserClient(conn)
}

func NewUsers(uc upbv1.UserClient) data.UserData {
	return &users{uc: uc}
}