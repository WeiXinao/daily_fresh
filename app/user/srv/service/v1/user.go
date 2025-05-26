package v1

import (
	"context"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	dv1 "github.com/WeiXinao/daily_your_go/app/user/srv/data/v1"
)

type UserSrv interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error) 
}

var _ UserSrv = (*userService)(nil)

type userService struct {
	userStore dv1.UserStore
}

func NewUserService(us dv1.UserStore) *userService {
	return &userService{
		userStore:	us, 
	}
}

type UserDTO struct {
	Name string
}

type UserDTOList struct {
	TotalCount int64 // 总数
	Items []*UserDTO  // 数据
}

/* 
	代码不方便写单元测试用例 
		1. data 层的接口必须先写好
		2. 我期望测试的时候测试底层的 data 层的数据按照我期望的返回
			1. 先手动去插入一些数据
			2. 去删除一些数据
		3. 如果 data 层的方法有 bug，坑，我们的代码想要具备好的可测试性
		
*/
func (u *userService) List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error) {
	// 这里是业务逻辑1
	doList, err := u.userStore.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 业务逻辑2
	// 代码不方便写单元测试用例

	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{Name: value.Name}	
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	// 业务逻辑3

	return &userDTOList, nil
}