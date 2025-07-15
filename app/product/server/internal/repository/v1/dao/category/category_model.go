package category

import "github.com/WeiXinao/daily_fresh/app/pkg/gormx"

type CategoryModel struct {
	gormx.BaseModel
	Name         string `gorm:"column:name;type:char(50);comment:'分类名称'"`
	ParentCid    int64  `gorm:"column:parent_cid;type:bigint;comment:'父分类id'"`
	CatLevel     int    `gorm:"column:cat_level;type:int;comment:comment:'层级'"`
	ShowStatus   int8   `gorm:"column:show_status;type:tinyint;comment:'是否显示[0-不显示,1-显示]'"`
	Sort         int    `gorm:"column:sort;type:int;comment:'排序'"`
	Icon         string `gorm:"column:icon;type:char(255);comment:'图标单位'"`
	ProductUint  string `gorm:"column:product_uint;type:char(50);comment:'计量单位'"`
	ProductCount int    `gorm:"column:product_count;type:int;comment:'商品单位'"`
}

func (*CategoryModel) TableName() string {
	return "pms_category"
}
