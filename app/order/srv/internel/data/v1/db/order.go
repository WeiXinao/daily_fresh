package db

import (
	"context"
	"strings"

	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/data/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	baseCode "github.com/WeiXinao/daily_your_go/gmicro/code"
	v1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"gorm.io/gorm"
)

type orders struct {
	db *gorm.DB
}

// Create 创建订单之后要删除对应的购物车记录
// Create implements data.OrderStore.
func (o *orders) Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
	if txn == nil {
		txn = o.db
	}
	err := txn.WithContext(ctx).Create(order).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Get implements data.OrderStore.
func (o *orders) Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error) {
	var order do.OrderInfoDO
	err := o.db.WithContext(ctx).Preload("OrderGoods").Where("order_sn", orderSn).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrOrderNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &order, nil
}

// List implements data.OrderStore.
func (o *orders) List(ctx context.Context, userID uint64, meta v1.ListMeta, orderBy []string) (*do.OrderInfoDOList, error) {
	var ret do.OrderInfoDOList
	// 分页
	var limit, offset int
	if meta.PageSize <= 0 {
		limit = 10
	} else {
		limit = meta.PageSize
	}
	if meta.Page > 0 {
		offset = (meta.Page - 1) * limit
	} 

	// 排序
	orderCond := strings.Join(orderBy, ",")

	query := o.db.Preload("OrderGoods")
	if offset > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if len(orderBy) > 0 {
		query = query.Order(orderCond)
	}

	err := query.WithContext(ctx).Find(&ret.Items).Count(&ret.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &ret, nil
}

// UpdateStatus implements data.OrderStore.
func (o *orders) Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
	if txn == nil {
		txn = o.db
	}
	err := txn.WithContext(ctx).Model(order).Updates(order).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

var _ data.OrderStore = (*orders)(nil)

func newOrders(factory *dataFactory) *orders {
	return &orders{
		db: factory.db,
	}
}
