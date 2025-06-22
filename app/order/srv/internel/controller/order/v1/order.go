package order

import (
	"context"

	opb "github.com/WeiXinao/daily_your_go/api/order/v1"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/do"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/domain/dto"
	"github.com/WeiXinao/daily_your_go/app/order/srv/internel/service/v1"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/xkit/slice"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ opb.OrderServer = (*OrderServer)(nil)

type OrderServer struct {
	opb.UnimplementedOrderServer
	srv service.ServiceFactory
}

// CartItemList implements proto.OrderServer.
func (o *OrderServer) CartItemList(context.Context, *opb.UserInfo) (*opb.CartItemListResponse, error) {
	panic("unimplemented")
}

// CreateCartItem implements proto.OrderServer.
func (o *OrderServer) CreateCartItem(context.Context, *opb.CartItemRequest) (*opb.ShopCartInfoResponse, error) {
	panic("unimplemented")
}

// 这个是给分布式事务 saga 调用的
// CreateOrder implements proto.OrderServer.
func (o *OrderServer) CreateOrder(ctx context.Context,req *opb.OrderRequest) (*emptypb.Empty, error) {
	orderInfo := &dto.OrderInfoDTO{
		OrderInfoDO: do.OrderInfoDO{
			User:         req.GetUserId(),
			OrderSn:      req.GetOrderSn(),
			Address:      req.GetAddress(),
			SignerName:   req.GetName(),
			SingerMobile: req.GetMobile(),
			Post:         req.GetPost(),
			OrderGoods: slice.Map(req.GetOrderItems(), func(idx int, src *opb.OrderItemReponse) *do.OrderGoodsDO {
				return &do.OrderGoodsDO{
					Goods: src.GetGoodsId(),
					GoodsName: src.GetGoodsName(),
					GoodsImage: src.GetGoodsImage(),
					GoodsPrice: src.GetGoodsPrice(),
					Nums: src.GetNums(),
				}
			}),
		},
	}
	err := o.srv.Order().Create(ctx, orderInfo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CreateOrderCom implements proto.OrderServer.
func (o *OrderServer) CreateOrderCom(context.Context, *opb.OrderRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// DeleteCartItem implements proto.OrderServer.
func (o *OrderServer) DeleteCartItem(context.Context, *opb.CartItemRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// OrderDetail implements proto.OrderServer.
func (o *OrderServer) OrderDetail(context.Context, *opb.OrderRequest) (*opb.OrderInfoDetailReponse, error) {
	panic("unimplemented")
}

// OrderList implements proto.OrderServer.
func (o *OrderServer) OrderList(context.Context, *opb.OrderFilterRequest) (*opb.OrderListResponse, error) {
	panic("unimplemented")
}

/*
订单提交的时候应该先生成订单号
订单号单独做一个接口，然后在订单提交的时候，会调用这个接口，生成订单号
订单查询以及一系列的关联， 我们应该采用 order_sn，不要再去采用 id 关联了
*/
// SubmitOrder implements proto.OrderServer.
func (o *OrderServer) SubmitOrder(ctx context.Context, req *opb.OrderRequest) (*emptypb.Empty, error) {
	err := o.srv.Order().Submit(ctx, &dto.OrderInfoDTO{
		OrderInfoDO: do.OrderInfoDO{
			User: req.GetUserId(),
			Address: req.GetAddress(),
			SignerName: req.GetName(),
			SingerMobile: req.GetMobile(),
			Post: req.GetPost(),
		},
	})
	if err != nil {
		log.Errorf("新建订单失败: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateCartItem implements proto.OrderServer.
func (o *OrderServer) UpdateCartItem(context.Context, *opb.CartItemRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateOrderStatus implements proto.OrderServer.
func (o *OrderServer) UpdateOrderStatus(context.Context, *opb.OrderStatus) (*emptypb.Empty, error) {
	panic("unimplemented")
}

func NewOrderServer(svc service.ServiceFactory) *OrderServer {
	return &OrderServer{
		srv: svc,
	}
}