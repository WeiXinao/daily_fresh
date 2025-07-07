package db

import (
	"context"
	"sync"
	"time"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	ipb "github.com/WeiXinao/daily_your_go/api/inventory/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/data/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"gorm.io/gorm"
)

var _ data.DataFactory = (*dataFactory)(nil)

type dataFactory struct {
	db          *gorm.DB
	invClient   ipb.InventoryClient
	goodsClient gpb.GoodsClient
}

// Begin implements data.DataFactory.
func (d *dataFactory) Begin() *gorm.DB {
	return d.db.Begin()
}

// Goods implements data.DataFactory.
func (d *dataFactory) Goods() gpb.GoodsClient {
	return d.goodsClient
}

// Inventorys implements data.DataFactory.
func (d *dataFactory) Inventorys() ipb.InventoryClient {
	return d.invClient
}

// Orders implements data.DataFactory.
func (d *dataFactory) Orders() data.OrderStore {
	return newOrders(d)
}

// ShopCarts implements data.DataFactory.
func (d *dataFactory) ShopCarts() data.ShopCartStore {
	return newShopCarts(d)
}

var (
	df   data.DataFactory
	once sync.Once
)

func GetDataFactory(mysqlOpts *options.MysqlOptions, registeryOpts *options.RegisteryOptions) (data.DataFactory, error) {
	db, err := getDB(mysqlOpts)
	if err != nil {
		return nil, err
	}
	registry, err := initServiceRegestry(registeryOpts)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	gc, err := connGoodsRpcServer(ctx, time.Hour, registry)
	if err != nil {
		return nil, err
	}
	ic, err := connInventoryRpcServer(ctx, time.Hour, registry)
	if err != nil {
		return nil, err
	}
	return &dataFactory{
		db: db,
		goodsClient: gc,
		invClient: ic,
	}, nil
}
