package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/domain/dto"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/go-redsync/redsync/v4"
)

const (
	inventoryLockPrefix = "inventory_"
	orderLockPrefix     = "order_"
)

type InventoryService interface {
	// 设置库存
	Create(ctx context.Context, inv *dto.InventoryDTO) error

	// 根据商品的 id 查询库存
	Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error)

	// 扣减库存
	Sell(ctx context.Context, ordersn string, detail []do.GoodsDetail) error

	// 归还库存
	Reback(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
}

type inventoryService struct {
	data data.DataFactory

	redsync *redsync.Redsync
}

// Create implements InventoryService.
func (i *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return i.data.Inventorys().Create(ctx, &inv.Inventory)
}

// Get implements InventoryService.
func (i *inventoryService) Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error) {
	inv, err := i.data.Inventorys().Get(ctx, goodsID)
	if err != nil {
		return nil, err
	}
	return &dto.InventoryDTO{
		Inventory: *inv,
	}, nil
}

// Reback implements InventoryService.
func (i *inventoryService) Reback(ctx context.Context, ordersn string, detail []do.GoodsDetail) error {
	log.Infof("订单%s归还库存", ordersn)

	txn := i.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Error("事务中出现异常，回滚事务")
		}
	}()

	// 库存归还的时候有不少细节
	// 1. 主动取消
	// 2. 网络问题引起的重试
	// 3. 超时取消 
	// 4. 退款取消
	mutex := i.redsync.NewMutex(fmt.Sprintf("%s:%s", orderLockPrefix, ordersn))
	if err := mutex.Lock(); err != nil {
		txn.Rollback()
		log.Errorf("订单%s获取锁失败", ordersn)
		return err
	}
	sellDetail, err := i.data.Inventorys().GetSellDetail(ctx, txn, ordersn)
	if err != nil {
		txn.Rollback()
		if _, err = mutex.Unlock(); err != nil {
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
		if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
			log.Errorf("订单%s没有找到库存扣减记录，忽略", ordersn)
			return nil
		}
		log.Errorf("订单%s获取库存扣减记录失败", ordersn)
		return err
	}

	if sellDetail.Status != 2 {
		log.Infof("订单%s库存记录已经归还，忽略", ordersn)
		if _, err = mutex.Unlock(); err != nil {
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
		return nil
	}

	for _, d := range detail {
		_, err := i.data.Inventorys().Get(ctx, uint64(d.Goods))
		if err != nil {
			txn.Rollback()
			if _, err = mutex.Unlock(); err != nil {
				log.Errorf("订单%s释放锁出现异常", ordersn)
			}
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}

		if err := i.data.Inventorys().Increase(ctx, txn, uint64(d.Goods), int(d.Num)); err != nil {
			txn.Rollback() // 回滚
			if _, err = mutex.Unlock(); err != nil {
				log.Errorf("订单%s释放锁出现异常", ordersn)
			}
			log.Errorf("订单%s归还库存失败", ordersn)
			return err
		}
	}

	err = i.data.Inventorys().UpdateStockSellDetailStatus(ctx, txn, ordersn, 2)
	if err != nil {
		txn.Rollback()
		if _, err = mutex.Unlock(); err != nil {
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
		log.Errorf("订单%s更新库存扣减记录状态失败", ordersn)
		return err
	}

	if _, err = mutex.Unlock(); err != nil {
		log.Errorf("订单%s释放锁出现异常", ordersn)
	}
	txn.Commit()
	return nil
}

// Sell implements InventoryService.
func (i *inventoryService) Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s扣减库存", ordersn)
	// 实际上批量扣减库存的时候，我们经常会先按照商品的 id 排序，然后从小到大逐个扣减库存，这样可以减少锁的竞争
	// 如果无序的话，那么就有可能订单 a，扣减 1，3，4，订单 B 扣减 3,2,1
	var detail = do.GoodsDetailList(details)
	sort.Sort(&detail)

	txn := i.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Error("事务中出现异常，回滚事务")
		}
	}()

	sellDetail := do.StockSellDetailDO{
		OrderSn: ordersn,
		Status:  1,
		Detail:  detail,
	}

	for _, d := range sellDetail.Detail {
		mutex := i.redsync.NewMutex(fmt.Sprintf("%s:%d", inventoryLockPrefix ,d.Goods))
		if err := mutex.Lock(); err != nil {
			log.Errorf("订单%s获取锁失败", ordersn)
		}

		inv, err := i.data.Inventorys().Get(ctx, uint64(d.Goods))
		if err != nil {
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}
		
		// 判断库存是否充足
		if inv.Stocks < d.Num {
			txn.Rollback() // 回滚事务
			log.Errorf("商品%d库存不足，需要%d，现有库存%d", d.Goods, d.Num, inv.Stocks)
			return errors.WithCode(code.ErrInvNotEnough, "库存不足")
		}
		// inv.Stocks -= d.Num

		err = i.data.Inventorys().Reduce(ctx, txn, uint64(d.Goods), int(d.Num))
		if err != nil {
			txn.Rollback()
			log.Errorf("订单%s扣减库存失败", ordersn)
			return err
		}

		// 释放锁
		if _, err = mutex.Unlock(); err != nil {
			txn.Rollback()
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
	}

	err := i.data.Inventorys().CreateStockSellDetail(ctx, txn, &sellDetail)
	if err != nil {
		txn.Rollback()
		log.Errorf("订单%s创建库存扣减记录失败", ordersn)
		return err
	} 

	txn.Commit()

	return nil
}

var _ InventoryService = &inventoryService{}

func newInventoryService(df data.DataFactory, redsync *redsync.Redsync) InventoryService {
	return &inventoryService{
		data:         df,
		redsync:      redsync,
	}
}
