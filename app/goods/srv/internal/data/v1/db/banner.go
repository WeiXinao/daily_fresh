package db

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	v1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

var _ data.BannerStore = (*banners)(nil)

type banners struct {
	db *gorm.DB
}

// Create implements data.BannerStore.
func (b *banners) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	panic("unimplemented")
}

// Delete implements data.BannerStore.
func (b *banners) Delete(ctx context.Context, id uint64) error {
	panic("unimplemented")
}

// List implements data.BannerStore.
func (b *banners) List(ctx context.Context, opts v1.ListMeta) (*do.BannerDOList, error) {
	panic("unimplemented")
}

// Update implements data.BannerStore.
func (b *banners) Update(ctx context.Context, txn *gorm.DB, bannner *do.BannerDO) error {
	panic("unimplemented")
}

func newBanners(factory *mysqlFactory) data.BannerStore {
	return &banners{
		db: factory.db,
	}
}
