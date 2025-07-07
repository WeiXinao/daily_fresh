package do

import "github.com/WeiXinao/daily_your_go/app/pkg/gormx"

type (
	ShoppingCartDO struct {
		gormx.BaseModel
		User    int32 `gorm:"type:int;index"` // 在购物车列表中我们需要查询当前用户的购物车记录
		Goods   int32 `gorm:"type:int;index"` // 加索引：我们需要查询的时候 1. 会影响插入性能 2. 会占用磁盘
		Nums    int32 `gorm:"type:int"`
		Checked bool  // 是否选中
	}

	ShoppingCartDOList struct {
		TotalCount int64 `json:"totalCount"`
		Items      []*ShoppingCartDO `json:"items"`
	}
)

func (ShoppingCartDO) TableName() string {
	return "shoppingcart"
}