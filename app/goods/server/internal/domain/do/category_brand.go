package do

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
)

type GoodsCategoryBrandDO struct {
	gormx.BaseModel
	CategoryID int32 `gorm:"type:int;not null;index:idx_category_brand,unique;comment:'分类ID'"`
	Category   CategoryDO
	BrandsID   int32 `gorm:"type:int;index:idx_category_brand,unique;comment:'品牌ID'"`
	Brands     BrandsDO
}

func (GoodsCategoryBrandDO) TableName() string {
	return "goodscategorybrand"
}

type (
	GoodsCategoryBrandDOList struct {
		TotalCount int64                   `json:"totalCount"`
		Items      []*GoodsCategoryBrandDO `json:"items"`
	}
)
