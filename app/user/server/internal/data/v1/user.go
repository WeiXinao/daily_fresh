package v1

import (
	"context"
	"time"

	"github.com/WeiXinao/daily_fresh/app/pkg/gormx"
	metav1 "github.com/WeiXinao/daily_fresh/pkg/common/meta/v1"
)

/*
	有数据访问的方法，一定要返回一个 error
	参数中最好有 ctx
*/
type UserStore interface {
	// 用户列表
	List(ctx context.Context, orderBy []string, opts metav1.ListMeta) (*UserDOList, error)

	// 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)

	// 通过用户 ID 查询用户
	GetByID(ctx context.Context, id uint64) (*UserDO, error)

	// 创建用户
	Create(ctx context.Context, user *UserDO) error

	// 更新用户
	Update(ctx context.Context, user *UserDO) error
}

/*
1. 密文 2. 密文不可反解
	1. 对称加密
	2. 非对称加密
	3. md5 信息摘要算法
	密码如果不可以反解，用户找回密码
*/
type UserDO struct {
	gormx.BaseModel
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女，male表示男'"`
	Role int `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}

func (u *UserDO) TableName() string {
	return "user"
}

type UserDOList struct {
	TotalCount int64 // 总数
	Items []*UserDO 
}
