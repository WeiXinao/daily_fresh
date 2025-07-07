package dto

import "github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"

type OrderInfoDTO struct {
	do.OrderInfoDO	
}

type OrderDTOList struct {
	TotalCount int64 `json:"totalCount"`
	Items []*OrderInfoDTO `json:"data"`
}