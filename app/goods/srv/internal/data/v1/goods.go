package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type GoodsStore interface {
	Get(ctx context.Context, id uint64) (*do.GoodsDO, error)
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error)
	ListByID(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *do.GoodsDO) error
	CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Delete(ctx context.Context, id uint64) error
	DeleteInTxn(ctx context.Context, txn *gorm.DB, id uint64) error

	Begin() *gorm.DB
}