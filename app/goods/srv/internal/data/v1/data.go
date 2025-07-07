package data

import "gorm.io/gorm"

type DataFactory interface {
	Goods() GoodsStore
	Categorys() CategoryStore
	Brands() BrandsStore
	Banner() BannerStore
	CategoryBrands() GoodsCategoryBrandStore

	Begin() *gorm.DB
}

