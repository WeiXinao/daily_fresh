package data

import (
	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	ipb "github.com/WeiXinao/daily_your_go/api/inventory/v1"
	"gorm.io/gorm"
)

type DataFactory interface {
	Orders() OrderStore
	ShopCarts() ShopCartStore
	Goods() gpb.GoodsClient
	Inventorys() ipb.InventoryClient

	Begin() *gorm.DB
}