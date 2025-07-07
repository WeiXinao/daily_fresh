package do

import (
	"time"

	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
)

type (
	OrderInfoDO struct {
		gormx.BaseModel

		OrderGoods []*OrderGoodsDO `gorm:"foreignKey:Order;references:ID"`

		User    int32  `gorm:"type:int;index"`
		OrderSn string `gorm:"type:varchar(30);index"` // 订单号，我们平台自己生成的订单号
		PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

		//status大家可以考虑使用iota来做
		Status     string  `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
		TradeNo    string  `gorm:"type:varchar(100) comment '交易号'"` // 交易号就是支付宝的订单号 查账
		OrderMount float32 // 订单的金额
		PayTime    *time.Time `gorm:"type:datetime"`

		Address      string `gorm:"type:varchar(100)"`
		SignerName   string `gorm:"type:varchar(20)"`
		SingerMobile string `gorm:"type:varchar(11)"`
		Post         string `gorm:"type:varchar(20)"` // 留言信息
	}

	OrderInfoDOList struct {
		TotalCount int64 `json:"totalCount"`
		Items      []*OrderInfoDO `json:"items"`
	}
)

func (OrderInfoDO) TableName() string {
	return "orderinfo"
}

type OrderGoodsDO struct {
	gormx.BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	// 把商品的信息保存下来，字段冗余，高并发系统中我们一般都不会遵循三范式
	// 做镜像，购买时的价格与现在的价格可能不一致
	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"` // 购买了多少件
}

func (OrderGoodsDO) TableName() string {
	return "ordergoods"
}