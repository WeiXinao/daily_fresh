package db

import (
	"context"
	"strings"
	"time"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	baseCode "github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"gorm.io/gorm"
)

var _ data.CategoryStore = (*categorys)(nil)

type categorys struct {
	db *gorm.DB
}

// Create implements data.CategoryStore.
func (c *categorys) Create(ctx context.Context, category *do.CategoryDO) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	err := c.db.Create(category).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Delete implements data.CategoryStore.
func (c *categorys) Delete(ctx context.Context, id uint64) error {
	err :=  c.db.Where("id = ?", id).Delete(&do.CategoryDO{}).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Get implements data.CategoryStore.
func (c *categorys) Get(ctx context.Context, ID uint64) (*do.CategoryDO, error) {
	category := &do.CategoryDO{}

	err := c.db.WithContext(ctx).Preload("SubCategory").Preload("SubCategory.SubCategory").First(category, ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(code.ErrCategoryNotFound, err.Error())
	}
	if err != nil {
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return category, nil
}

// ListAll implements data.CategoryStore.
func (c *categorys) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	ret := &do.CategoryDOList{}

	// 排序
	order := strings.Join(orderby, ",")

	var err error
	query := c.db.Where("level=1").Preload("SubCategory.SubCategory")
	if len(orderby) > 0 {
		err = query.Order(order).Find(&ret.Items).Error
	} else {
		err = query.Find(&ret.Items).Error
	}
	if err != nil {
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return ret, nil
}

// Update implements data.CategoryStore.
func (c *categorys) Update(ctx context.Context, category *do.CategoryDO) error {
	category.UpdatedAt = time.Now()
	err := c.db.Model(&do.CategoryDO{}).Updates(category).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// func NewCategory(db *gorm.DB) *categorys {
// 	return &categorys{
// 		db: db,
// 	}
// }

func newCategory(factory *mysqlFactory) *categorys {
	return &categorys{
		db: factory.db,
	}
}