package db

import (
	"context"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	v1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

var _ data.GoodsCategoryBrandStore = (*categoryBrands)(nil)

type categoryBrands struct {
	db *gorm.DB
}

// Create implements data.GoodsCategoryBrandStore.
func (c *categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	panic("unimplemented")
}

// Delete implements data.GoodsCategoryBrandStore.
func (c *categoryBrands) Delete(ctx context.Context, id uint64) error {
	panic("unimplemented")
}

// List implements data.GoodsCategoryBrandStore.
func (c *categoryBrands) List(ctx context.Context, opts v1.ListMeta, orderby []string) (*do.GoodsCategoryBrandDOList, error) {
	panic("unimplemented")
}

// Update implements data.GoodsCategoryBrandStore.
func (c *categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	panic("unimplemented")
}
