package do

import "github.com/WeiXinao/daily_your_go/app/pkg/gormx"

type GoodsSearchDO struct {
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	BrandsID   int32 `json:"brand_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	ShopPrice   float32 `json:"shop_price"`
	GoodsBrief  string  `json:"goods_brief"`
}

func (GoodsSearchDO) GetIndexName() string {
	return "goods"	
}

func (GoodsSearchDO) GetMapping() string {
	goodsMapping := `
	{
		"mappings" : {
			"properties" : {
				"brands_id" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"fav_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"ship_free" : {
					"type" : "boolean"
				},
				"shop_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
	return goodsMapping
}

type GoodsSearchDOList struct {
	TotalCount int64 `json:"totalCount,omitempty"`
	Items []*GoodsSearchDO `json:"items"`
}

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
