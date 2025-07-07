package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/data/v1"
	"github.com/WeiXinao/daily_your_go/app/inventory/srv/internal/domain/do"
	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlStore struct {
	db *gorm.DB
}

func (m *mysqlStore) Inventorys() data.InventoryStore {
	return newInventorys(m)
}

var _ data.DataFactory = &mysqlStore{}

var (
	dbFactory data.DataFactory
	once      sync.Once
)

// 对于复杂的初始化过程，使用工厂模式
func GetDBFactoryOr(mysqlOpts *options.MysqlOptions) (data.DataFactory, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)

		//希望大家自己可以去封装logger
		newLogger := gormx.New(
			gormx.Config{
				SlowThreshold:             time.Second,                         // 慢 SQL 阈值
				LogLevel:                  logger.LogLevel(mysqlOpts.LogLevel), // 日志级别
				IgnoreRecordNotFoundError: true,                                // 忽略ErrRecordNotFound（记录未找到）错误
			},
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return
		}

		sqlDB, _ := db.DB()
		dbFactory = &mysqlStore{
			db: db,
		}

		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifeTime)
	})

	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func MigrateDatabase(db *gorm.DB) error {
	//if err := db.AutoMigrate(&v12.Inventory{}); err != nil {
	//	return errors.Wrap(err, "migrate inventory model failed")
	//}

	if err := db.AutoMigrate(&do.StockSellDetailDO{}); err != nil {
		return errors.Wrap(err, "migrate brand model failed")
	}

	return nil
}

func (ds *mysqlStore) Begin() *gorm.DB {
	return ds.db.Begin()
}
