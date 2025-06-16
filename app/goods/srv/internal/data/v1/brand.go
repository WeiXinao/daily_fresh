package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)


type	BrandsStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error)
	Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error
	Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error
	Delete(ctx context.Context, id uint64) error
}