package do

import "github.com/WeiXinao/daily_your_go/app/pkg/gormx"

type GoodsDO struct {
	gormx.BaseModel

	CategoryID int32 `gorm:"type:int;not null;comment:'分类ID'"`
	Category CategoryDO 
	BrandsID int32 `gorm:"type:int;not null;comment:'品牌ID'"`
	Brands BrandsDO 

	OnSale bool `gorm:"default:false;not null;comment:'是否上架'"`
	ShipFree bool `gorm:"default:false;not null;comment:'是否免运费'"`
	IsNew bool `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot bool `gorm:"default:false;not null;comment:'是否热销'"`

	Name string `gorm:"type:varchar(50);not null;comment:'品牌名称'"`
	GoodsSn string `gorm:"type:varchar(50);not null;comment:'商品编号'"`
	ClickNum int32 `gorm:"type:int;default:0;not null;comment:'点击量'"`
	SoldNum int32 `gorm:"type:int;default:0;not null;comment:'销量'"`
	FavNum int32 `gorm:"type:int;default:0;not null;comment:'收藏量'"`
	MarketPrice float32 `gorm:"not null;comment:'市场价'"`
	ShopPrice float32 `gorm:"not null;comment:'本店价'"`
	GoodsBrief string `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	Images gormx.GormList `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages gormx.GormList `gorm:"type:varchar(1000);not null;comment:'商品描述图片'"`
	GoodsFrontImage string `gorm:"type:varchar(200);not null;comment:'商品封面图片'"`
}

type GoodsDOList struct {
	TotalCount int64 `json:"totalCount"`
	Items []*GoodsDO `json:"items"`
}
