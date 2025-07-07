package dto

import "github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"

type ShopCartDTO struct {
	do.ShoppingCartDO
}

type ShopCartDTOList struct {
	TotalCount int64       `json:"totalCount"`
	Items      []*ShopCartDTO `json:"items"`
}
