package db

import (
	"context"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	v1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

var _ data.BrandsStore = (*brands)(nil)

type brands struct {
	db *gorm.DB
}

// Create implements data.BrandsStore.
func (b *brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	panic("unimplemented")
}

// Delete implements data.BrandsStore.
func (b *brands) Delete(ctx context.Context, id uint64) error {
	panic("unimplemented")
}

// List implements data.BrandsStore.
func (b *brands) List(ctx context.Context, opts v1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	panic("unimplemented")
}

// Update implements data.BrandsStore.
func (b *brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	panic("unimplemented")
}
