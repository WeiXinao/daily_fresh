package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	"gorm.io/gorm"
)

type GoodsStore interface {
	Get(ctx context.Context, id uint64) (*do.GoodsDO, error)
	ListByID(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	Create(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Update(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Delete(ctx context.Context, id uint64) error
}