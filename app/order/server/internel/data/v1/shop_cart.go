package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type ShopCartStore interface {
	List(ctx context.Context, userID uint64, checked bool, meta metav1.ListMeta, orderBy []string) (*do.ShoppingCartDOList, error)
	Create(ctx context.Context, shopCartItem *do.ShoppingCartDO) error
	Get(ctx context.Context, userID uint64, goodsID uint64) (*do.ShoppingCartDO, error)
	UpdateNum(ctx context.Context, shopCartItem *do.ShoppingCartDO) error
	Delete(ctx context.Context, id uint64) error
	ClearCheck(ctx context.Context, userID uint64) error
	DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []uint64) error
}