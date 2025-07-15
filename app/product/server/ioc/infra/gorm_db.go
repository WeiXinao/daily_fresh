package infra

import (
	"fmt"

	"github.com/WeiXinao/daily_fresh/app/pkg/gormx"
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/WeiXinao/daily_fresh/pkg/gmicro/code"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(mysqlOpts *options.MysqlOptions) (*gorm.DB, error) {
	logger := gormx.New(gormx.Config{
		LogLevel: logger.LogLevel(mysqlOpts.LogLevel),
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlOpts.Username, mysqlOpts.Password, mysqlOpts.Host, mysqlOpts.Port, mysqlOpts.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, err.Error())
	}

	sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifeTime)

	return db, nil
}