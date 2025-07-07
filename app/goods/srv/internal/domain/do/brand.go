package do

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
)

type BrandsDO struct {
	gormx.BaseModel
	Name string `gorm:"type:varchar(20);not null;uniqueIndex:uk_brand_name;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌logo'"`
}

func (BrandsDO) TableName() string {
	return "brands"
}

type BrandsDOList struct {
		TotalCount int64 `json:"totalCount"`
		Items []*BrandsDO `json:"items"`
	}

