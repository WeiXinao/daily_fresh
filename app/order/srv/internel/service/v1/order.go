package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	gpb "github.com/WeiXinao/daily_your_go/api/goods/v1"
	ipb "github.com/WeiXinao/daily_your_go/api/inventory/v1"
	opb "github.com/WeiXinao/daily_your_go/api/order/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/data/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/dto"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	metav1 "github.com/WeiXinao/daily_your_go/pkg/common/meta/v1"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/xkit/slice"
	"github.com/dtm-labs/dtmgrpc"
)

type OrderSrv interface {
	Get(ctx context.Context, orderSn string) (*dto.OrderInfoDTO, error)
	List(ctx context.Context, userId uint64, meta metav1.ListMeta, orderBy []string) (*dto.OrderDTOList, error)
	Submit(ctx context.Context, order *dto.OrderInfoDTO) error
	Create(ctx context.Context, order *dto.OrderInfoDTO) error // 这是 create 的补偿方法
	CreateCom(ctx context.Context, order *dto.OrderInfoDTO) error
	Update(ctx context.Context, order *dto.OrderInfoDTO) error
}

var _ OrderSrv = (*orderService)(nil)

type orderService struct {
	data    data.DataFactory
	dtmOpts *options.DtmOptions
}

// CreateCom implements OrderSrv.
func (o *orderService) CreateCom(ctx context.Context, order *dto.OrderInfoDTO) error {
	/*
		1. 删除 orderinfo 表
		2. 删除 ordergoods 表
		3. 根据 order 找到对应的购物车条目，恢复购物车条目
	*/
	// 其实不用回滚
	// 你应该先查询订单是否已经存在，如果已经存在删除相关记录即可，同时恢复购物车记录
	return nil
}

// Create implements OrderSrv.
func (o *orderService) Create(ctx context.Context, order *dto.OrderInfoDTO) error {
	/*
		1. 生成 orderinfo 表
		2. 生成 ordergoods 表
		3. 根据 order 找到对应的购物车条目，删除购物车条目
	*/
	ids := slice.Map(order.OrderGoods, func(idx int, src *do.OrderGoodsDO) int32 {
		return src.Goods
	})
	goods, err := o.data.Goods().BatchGetGoods(ctx, &gpb.BatchGoodsIdInfo{
		Id: 	ids,
	})
	if err != nil {
		log.Errorf("批量获取商品信息失败，gooodsIds: %v, err: %w", ids, err)
		return errors.FromGrpcError(err)
	}

	id2Goods := slice.ToMap(goods.Data, func(element *gpb.GoodsInfoResponse) int32 {
		return element.Id
	})

	var orderAmount float32
	for _, value := range order.OrderGoods {
		good := id2Goods[value.Goods]
		value.GoodsName = good.GetName()
		value.GoodsPrice = good.GetShopPrice()
		value.GoodsImage = good.GetGoodsFrontImage()
		orderAmount += good.GetShopPrice() * float32(value.Nums)
	}
	order.OrderMount = orderAmount

	txn := o.data.Begin()
	defer func()  {
		if e := recover(); e != nil {
			txn.Rollback()
			log.Error("新建订单事务进行中出现异常，回滚")
			return
		}
	}()

	err = o.data.Orders().Create(ctx, txn, &order.OrderInfoDO)
	if err != nil {
		txn.Rollback()
		log.Errorf("创建订单失败, err: %w", err)
		return err
	}

	goodsIds := slice.Map(ids, func(idx int, src int32) uint64 { return uint64(src) })
	err = o.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, uint64(order.User), goodsIds)
	if err != nil {
		txn.Rollback()
		log.Errorf("删除购物车失败, err: %v", err)
		return err
	}

	txn.Commit()
	return nil
}

// 订单号的生成，订单号-雪花算法，目前的订单号生成算法有问题：不是递增的
func generateOrderSn(userId int32) string {
	/*
	 订单号的生成规则：
	 年月日时分秒 + 用户id + 2 位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

// Submit implements OrderSrv.
func (o *orderService) Submit(ctx context.Context, order *dto.OrderInfoDTO) error {
	// 先从购物车中获取商品
	list, err := o.data.ShopCarts().List(ctx, uint64(order.User), true, metav1.ListMeta{}, []string{})
	if err != nil {
		log.Errorf("获取购物车信息失败，err:%w", err)
		return err
	}
	
	if len(list.Items) == 0 {
		log.Error("购物车中没有商品，无法下单")
		return errors.WithCode(code.ErrNoGoodsSelected, "没有选择商品")
	}

	order.OrderGoods = slice.Map(list.Items, func(idx int, src *do.ShoppingCartDO) *do.OrderGoodsDO {
		return &do.OrderGoodsDO{
			Goods: src.Goods,
			Nums: src.Nums,
		}
	})

	// 基于可靠消息最终一致性的思想，saga 事务来解决订单生成的问题
	// orderSn := generateOrderSn(order.User)
	// order.OrderSn = orderSn
	req := &ipb.SellInfo{
		OrderSn: order.OrderSn,
		GoodsInfo: slice.Map(order.OrderGoods, func(idx int, src *do.OrderGoodsDO) *ipb.GoodsInvInfo {
			return &ipb.GoodsInvInfo{
				GoodId: src.Goods,
				Num:    src.Nums,
			}
		}),
	}
	var (
		invRpcUrl = "discovery:///daily-your-go-inventory-srv"
		orderRpcUrl = "discovery:///daily-your-go-order-srv"
	)
	saga := dtmgrpc.NewSagaGrpc(o.dtmOpts.GrpcServer, order.OrderSn).
		Add(invRpcUrl+ipb.Inventory_Sell_FullMethodName, invRpcUrl+ipb.Inventory_Reback_FullMethodName, req).
		Add(orderRpcUrl+opb.Order_CreateOrder_FullMethodName, orderRpcUrl+opb.Order_CreateOrderCom_FullMethodName, req)
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		return errors.WithCode(code.ErrSubmitOrder, "提交订单失败")
	}
	// 通过 OrderSn查询一下，当前转态如果一直 Submit 那么你就一直不要给前端返回，
	// 如果是 failed 那么你提示给前端说下单失败，重新下单
	return nil
}

// Get implements OrderSrv.
func (o *orderService) Get(ctx context.Context, orderSn string) (*dto.OrderInfoDTO, error) {
	order, err := o.data.Orders().Get(ctx, orderSn)
	if err != nil {
		return nil, err
	}
	return &dto.OrderInfoDTO{
		OrderInfoDO: *order,
	}, nil
}

// List implements OrderSrv.
func (o *orderService) List(ctx context.Context, userId uint64, meta metav1.ListMeta, orderBy []string) (*dto.OrderDTOList, error) {
	orders, err := o.data.Orders().List(ctx, userId, meta, orderBy)
	if err != nil {
		return nil, err
	}
	return &dto.OrderDTOList{
		TotalCount: orders.TotalCount,
		Items: slice.Map(orders.Items, func(idx int, src *do.OrderInfoDO) *dto.OrderInfoDTO {
			return &dto.OrderInfoDTO{
				OrderInfoDO: *src,
			}
		}),
	}, nil
}

// Update implements OrderSrv.
func (o *orderService) Update(ctx context.Context, order *dto.OrderInfoDTO) error {
	panic("unimplemented")
}

func newOrderService(data data.DataFactory, dtmOpts *options.DtmOptions) OrderSrv {
	return &orderService{
		data:    data,
		dtmOpts: dtmOpts,
	}
}
