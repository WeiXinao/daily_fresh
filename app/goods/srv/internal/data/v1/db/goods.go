package db

import (
	"context"
	"strings"
	"time"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	baseCode "github.com/WeiXinao/daily_your_go/gmicro/code"
	v1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"gorm.io/gorm"
)

var _ data.GoodsStore = (*goods)(nil)

type goods struct {
	db *gorm.DB
}

// List implements data.GoodsStore.
func (g *goods) List(ctx context.Context, orderby []string, opts v1.ListMeta) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}
	// 分页
	var limit, offset int
	if opts.PageSize <= 0 {
		limit = 10
	}
	
	if opts.Page <= 0 {
		opts.Page = 1
	}
	offset = (opts.Page - 1) * limit

	// 排序
	order := strings.Join(orderby, ",")

	// 查询
	var err error
	query := g.db.Model(&do.GoodsDO{}).Preload("Category").Preload("Brands").Count(&ret.TotalCount)
	if len(strings.TrimSpace(order)) == 0 {
		err = query.Limit(limit).Offset(offset).Find(&ret.Items).Error
	} else {
		err = query.Order(order).Limit(limit).Offset(offset).Find(&ret.Items).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrGoodsNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return ret, nil
}

// Create implements data.GoodsStore.
func (g *goods) Create(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	goods.CreatedAt = time.Now()
	goods.UpdatedAt = time.Now()
	err := g.db.Create(goods).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Delete implements data.GoodsStore.
func (g *goods) Delete(ctx context.Context, id uint64) error {
	err :=  g.db.Where("id = ?", id).Delete(&do.GoodsDO{}).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}

// Get implements data.GoodsStore.
func (g *goods) Get(ctx context.Context, id uint64) (*do.GoodsDO, error) {
	goods := &do.GoodsDO{}
	err := g.db.Preload("Category").Preload("Brands").First(goods, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(code.ErrGoodsNotFound, err.Error())
	}
	if err != nil {
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return goods, nil
}

// ListByID implements data.GoodsStore.
func (g *goods) ListByID(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}

	// 排序
	order := strings.Join(orderby, ",")

	// 查询
	var err error
	query := g.db.Model(&do.GoodsDO{}).Preload("Category").Preload("Brands").Count(&ret.TotalCount)
	if len(strings.TrimSpace(order)) == 0 {
		err = query.Where("id in (?)", ids).Find(&ret.Items).Error
	} else {
		err = query.Where("id in (?)", ids).Order(order).Find(&ret.Items).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrGoodsNotFound, err.Error())
		}
		return nil, errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return ret, nil
}

// Update implements data.GoodsStore.
func (g *goods) Update(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	goods.UpdatedAt = time.Now()
	err := g.db.Model(&do.GoodsDO{}).Updates(goods).Error
	if err != nil {
		return errors.WithCode(baseCode.ErrDatabase, err.Error())
	}
	return nil
}
