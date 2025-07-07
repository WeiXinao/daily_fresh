package data

import "gorm.io/gorm"

type DataFactory interface {
	Inventorys() InventoryStore

	Begin() *gorm.DB
}
