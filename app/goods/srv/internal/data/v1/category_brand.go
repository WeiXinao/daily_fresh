package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type GoodsCategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandDOList, error)
	Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Delete(ctx context.Context, id uint64) error
}
