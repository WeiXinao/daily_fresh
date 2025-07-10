package db

import (
	"context"
	"strings"
	"time"

	"github.com/WeiXinao/daily_fresh/app/pkg/code"
	udv1 "github.com/WeiXinao/daily_fresh/app/user/server/internal/data/v1"
	baseCode "github.com/WeiXinao/daily_fresh/pkg/gmicro/code"
	metav1 "github.com/WeiXinao/daily_fresh/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"gorm.io/gorm"
)

var _ udv1.UserStore = (*users)(nil)

type users struct {
	db *gorm.DB
}

// Create 
// 	@description 创建用户
//	@param ctx 
//	@param user 
//	@return error 
func (u *users) Create(ctx context.Context, user *udv1.UserDO) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := u.db.Create(user).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// GetByID 
//  @description 根据用户 ID 获取用户信息
//	@param ctx 
//	@param id：用户 ID
//	@return *udv1.UserDO 
//	@return error 
func (u *users) GetByID(ctx context.Context, id uint64) (*udv1.UserDO, error) {
	user := udv1.UserDO{}
	err := u.db.First(&user, id).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &user, nil
}

// GetByMobile 
//  @description 根据手机号获取用户信息
//	@param ctx 
//	@param mobile 
//	@return *udv1.UserDO 
//	@return error 
func (u *users) GetByMobile(ctx context.Context, mobile string) (*udv1.UserDO, error) {
	user := udv1.UserDO{}

	// err 是 gorm 的 error，这种 error 我们经量不要抛出去
	err := u.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &user, nil
}

// Update 
// 	@descriptin 更新用户信息
//	@param ctx 
//	@param user 
//	@return error 
func (u *users) Update(ctx context.Context, user *udv1.UserDO) error {
	user.UpdatedAt = time.Now()
	err := u.db.Model(&udv1.UserDO{}).Updates(user).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// List 
// 	@descriptin 分页查询用户列表，凡是列表页返回的时候都应该返回总共有多少个
//	@param ctx 
//	@param opts 
//	@return *udv1.UserDOList 
//	@return error 
func (u *users) List(ctx context.Context, orderBy []string, opts metav1.ListMeta) (*udv1.UserDOList, error) {
	ret := &udv1.UserDOList{}
	// 分页
	var limit, offset int
	if opts.PageSize <= 0 {
		limit = 10
	}
	
	if opts.Page <= 0 {
		opts.Page = 1
	}
	offset = (opts.Page - 1) * limit

	// 排序
	order := strings.Join(orderBy, ",")

	// 查询
	var err error
	if len(strings.TrimSpace(order)) == 0 {
		err = u.db.Model(&udv1.UserDO{}).Count(&ret.TotalCount).
			Limit(limit).Offset(offset).Find(&ret.Items).Error
	} else {
		err = u.db.Model(&udv1.UserDO{}).Count(&ret.TotalCount).
			Order(order).Limit(limit).Offset(offset).Find(&ret.Items).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return ret, nil
}

func NewUsers(db *gorm.DB) udv1.UserStore {
	return &users{db: db}
}