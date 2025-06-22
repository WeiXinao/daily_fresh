package controller

import (
	"context"

	invpb "github.com/WeiXinao/daily_your_go/api/inventory/v1"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/domain/dto"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/service/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type inventoryServer struct {
	invpb.UnimplementedInventoryServer
	srv service.ServiceFactory
}

// 设置库存
func (is *inventoryServer) SetInv(ctx context.Context, info *invpb.GoodsInvInfo) (*emptypb.Empty, error) {
	invDTO := &dto.InventoryDTO{}
	invDTO.Inventory.Goods = info.GoodId
	invDTO.Inventory.Stocks = info.Num
	err := is.srv.Inventory().Create(ctx, invDTO)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (is *inventoryServer) InvDetail(ctx context.Context, info *invpb.GoodsInvInfo) (*invpb.GoodsInvInfo, error) {
	inv, err := is.srv.Inventory().Get(ctx, uint64(info.GoodId))
	if err != nil {
		return nil, err
	}
	return &invpb.GoodsInvInfo{
		GoodId: inv.Inventory.Goods,
		Num:    inv.Inventory.Stocks,
	}, nil
}

func (is *inventoryServer) Sell(ctx context.Context, info *invpb.SellInfo) (*emptypb.Empty, error) {
	var detail []do.GoodsDetail
	for _, value := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{Goods: value.GoodId, Num: value.Num})
	}
	err := is.srv.Inventory().Sell(ctx, info.OrderSn, detail)
	if err != nil {
		if errors.IsCode(err, code.ErrInvNotEnough) {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		return nil, err
	}
	//time.Sleep(5 * time.Second)
	// return nil, status.Errorf(codes.Aborted, "")
	return &emptypb.Empty{}, nil
}

func (is *inventoryServer) Reback(ctx context.Context, info *invpb.SellInfo) (*emptypb.Empty, error) {
	log.Infof("订单%s归还库存", info.OrderSn)
	var detail []do.GoodsDetail
	for _, v := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{Goods: v.GoodId, Num: v.Num})
	}
	err := is.srv.Inventory().Reback(ctx, info.OrderSn, detail)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func NewInventoryServer(srv service.ServiceFactory) *inventoryServer {
	return &inventoryServer{srv: srv}
}

var (
	_ invpb.InventoryServer = &inventoryServer{}
)
