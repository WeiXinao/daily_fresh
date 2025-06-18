package data

import gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"

type DataFactory interface {
	Goods() gpb.GoodsClient
	User() UserData
}
