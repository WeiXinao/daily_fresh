package db

import (
	"fmt"
	"sync"

	"github.com/WeiXinao/daily_your_go/app/pkg/gormx"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// 这个方法会返回 gorm 连接
// 这个方法应该返回的是全局的一个变量，如果一开始的时候没有初始化好，那么就初始化一次，
// 后续直接拿到这个变量就可以了
func getDB(mysqlOpts *options.MysqlOptions) (*gorm.DB, error) {
	if mysqlOpts == nil && db == nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}

	var (
		err error
	)
	dbOnce.Do(func() {
		logger := gormx.New(gormx.Config{
			LogLevel: logger.LogLevel(mysqlOpts.LogLevel),
		})

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlOpts.Username, mysqlOpts.Password, mysqlOpts.Host, mysqlOpts.Port, mysqlOpts.Database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger,
		})
		if err != nil {
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			return
		}

		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifeTime)
	})
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory, err: %w", err)
	}
	return db, nil
}
