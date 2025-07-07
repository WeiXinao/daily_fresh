package do

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
)


type CategoryDO struct {
	gormx.BaseModel        
	Name             string      `gorm:"type:varchar(20);not null;comment:'分类名称'" json:"name"`
	ParentCategoryID int32       `gorm:"type:int;comment:'父分类ID'" json:"parent"`
	ParentCategory   *CategoryDO   `json:"-"`
	SubCategory      []*CategoryDO `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1;comment:'分类级别'" json:"level"`
	IsTab            bool        `gorm:"default:false;not null;comment:'是否显示在Tab栏'" json:"is_tab"`
}

func (CategoryDO) TableName() string {
	return "category"
}

type CategoryDOList struct {
		TotalCount int64 `json:"totalCount"`
		Items []*CategoryDO `json:"items"`
	}
