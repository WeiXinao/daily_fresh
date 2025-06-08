package v1

import (
	"context"
	"crypto/sha512"
	"fmt"

	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	dv1 "github.com/WeiXinao/daily_your_go/app/user/srv/data/v1"
	v1 "github.com/WeiXinao/daily_your_go/app/user/srv/data/v1"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/anaskhan96/go-password-encoder"
)

type UserSrv interface {
	List(ctx context.Context, orderBy []string, opts metav1.ListMeta) (*UserDTOList, error)
	Create(ctx context.Context, user *UserDTO) error
	Update(ctx context.Context, user *UserDTO) error
	GetByID(ctx context.Context, id uint64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
}

var _ UserSrv = (*userService)(nil)

type userService struct {
	userStore dv1.UserStore
}

// Create implements UserSrv.
func (u *userService) Create(ctx context.Context, user *UserDTO) error {
	// 先判断用户是否存在
	_, err := u.userStore.GetByMobile(ctx, user.Mobile)
	if err == nil {
		return errors.WithCode(code.ErrUserAlreadyExists, "手机号重复，用户已存在")
	}
	if !errors.IsCode(err, code.ErrUserNotFound) {
		return err
	}

	// 密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(user.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	return u.userStore.Create(ctx, &user.UserDO)
}

// GetByID implements UserSrv.
func (u *userService) GetByID(ctx context.Context, id uint64) (*UserDTO, error) {
	userDO, err := u.userStore.GetByID(ctx, id)
	if err!= nil {
		return nil, err
	}
	return &UserDTO{*userDO}, nil
}

// GetByMobile implements UserSrv.
func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	udo, err := u.userStore.GetByMobile(ctx, mobile)
	if err!= nil {
		return nil, err	
	}
	return &UserDTO{
		UserDO: *udo,
	}, nil
}

// Update implements UserSrv.
func (u *userService) Update(ctx context.Context, user *UserDTO) error {
	// 先查询用户是否存在
	_, err := u.userStore.GetByID(ctx, uint64(user.ID))
	if err != nil {
		return err
	}
	return u.userStore.Update(ctx, &user.UserDO)
}

/*
代码不方便写单元测试用例
 1. data 层的接口必须先写好
 2. 我期望测试的时候测试底层的 data 层的数据按照我期望的返回
 1. 先手动去插入一些数据
 2. 去删除一些数据
 3. 如果 data 层的方法有 bug，坑，我们的代码想要具备好的可测试性
*/
func (u *userService) List(ctx context.Context, orderBy []string, opts metav1.ListMeta) (*UserDTOList, error) {
	// 这里是业务逻辑1
	doList, err := u.userStore.List(ctx, orderBy, opts)
	if err != nil {
		return nil, err
	}

	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{*value}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	return &userDTOList, nil
}

func NewUserService(us dv1.UserStore) *userService {
	return &userService{
		userStore: us,
	}
}

type UserDTO struct {
	v1.UserDO
}

type UserDTOList struct {
	TotalCount int64      // 总数
	Items      []*UserDTO // 数据
}
