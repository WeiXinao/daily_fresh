package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type BannerStore interface {
	List(ctx context.Context, opts metav1.ListMeta) (*do.BannerDOList, error)
	Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error
	Update(ctx context.Context, txn *gorm.DB, bannner *do.BannerDO) error
	Delete(ctx context.Context, id uint64) error
}
