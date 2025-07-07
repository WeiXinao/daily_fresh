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

var _ data.ShopCartStore = (*shopCarts)(nil)

type shopCarts struct {
	db *gorm.DB
}

// 这个在事务中执行，建议大家使用消息队列来实现
// DeleteByGoodsIDs implements data.ShopCartStore.
func (s *shopCarts) DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []uint64) error {
	if txn == nil {
		txn = s.db
	}
	err := txn.WithContext(ctx).Where("user = ? AND goods in (?)", userID, goodsIDs).Delete(&do.ShoppingCartDO{}).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// 删除选中商品的购物车记录，下订单了
// 从架构上来讲，这种实现有两种方式
// 下单后，直接执行购物车的记录，比较简单
// 下单后什么都不做，直接给 rocketmq 发送一个消息，然后由 rocketmq 来执行删除购物车记录

// 清空 check 状态
// ClearCheck implements data.ShopCartStore.
func (s *shopCarts) ClearCheck(ctx context.Context, userID uint64) error {
	panic("unimplemented")
}

// Create implements data.ShopCartStore.
func (s *shopCarts) Create(ctx context.Context, shopCartItem *do.ShoppingCartDO) error {
	if err := s.db.WithContext(ctx).Create(shopCartItem).Error; err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Delete implements data.ShopCartStore.
func (s *shopCarts) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&do.ShoppingCartDO{}).Error
}

// Get implements data.ShopCartStore.
func (s *shopCarts) Get(ctx context.Context, userID uint64, goodsID uint64) (*do.ShoppingCartDO, error) {
	var shopCart do.ShoppingCartDO
	err := s.db.WithContext(ctx).Where("user = ? and goods = ?", userID, goodsID).First(&shopCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrShopCartItemNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &shopCart, nil
}

// List implements data.ShopCartStore.
func (s *shopCarts) List(ctx context.Context, userID uint64, checked bool, meta v1.ListMeta, orderBy []string) (*do.ShoppingCartDOList, error) {
	var ret do.ShoppingCartDOList
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

	query := s.db
	if offset > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if len(orderBy) > 0 {
		query = query.Order(orderCond)
	}
	if userID > 0 {
		query = query.Where("user = ?", userID)
	}

	// 加入购物车状态
	if checked {
		query = query.Where("checked = ?", true)
	}

	err := query.WithContext(ctx).Find(&ret.Items).Count(&ret.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return &ret, nil
}

// UpdateNum implements data.ShopCartStore.
func (s *shopCarts) UpdateNum(ctx context.Context, shopCartItem *do.ShoppingCartDO) error {
	return s.db.WithContext(ctx).Model(do.ShoppingCartDO{}).
		Where("user = ? AND goods = ?", shopCartItem.User, shopCartItem.Goods).
		Update("nums", shopCartItem.Nums).
		Update("checked", shopCartItem.Checked).
		Error
}

func newShopCarts(factory *dataFactory) *shopCarts {
	return &shopCarts{
		db: factory.db,
	}
}
