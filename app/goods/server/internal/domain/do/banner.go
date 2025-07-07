package do

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
)

type BannerDO struct {
	gormx.BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:'轮播图'"`
	Url   string `gorm:"type:varchar(200);not null;comment:'跳转链接'"`
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图顺序'"`
}

func (BannerDO) TableName() string {
	return "banner"
}

type BannerDOList struct {
	TotalCount int64       `json:"totalCount"`
	Items      []*BannerDO `json:"items"`
}
