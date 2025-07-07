package data

import (
	"context"

	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type OrderStore interface {
	Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error)
	List(ctx context.Context, userID uint64, meta metav1.ListMeta, orderBy []string) (*do.OrderInfoDOList, error)
	Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error
	Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error
}