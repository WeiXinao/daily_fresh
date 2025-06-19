package service

import (
	"fmt"

	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type ServiceFactory interface {
	Inventory() InventoryService
}

var _ ServiceFactory = &service{}

type service struct {
	data         data.DataFactory
	redisOptions *options.RedisOptions
	redsync *redsync.Redsync
}

// InventoryService implements ServiceFactory.
func (s *service) Inventory() InventoryService {
	panic("unimplemented")
}

func NewService(df data.DataFactory, redisOptions *options.RedisOptions) ServiceFactory {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", redisOptions.Host, redisOptions.Port),
	})
	pool := goredis.NewPool(client) 
	redsync := redsync.New(pool)
	return &service{
		data:         df,
		redisOptions: redisOptions,
		redsync:      redsync,
	}
}

