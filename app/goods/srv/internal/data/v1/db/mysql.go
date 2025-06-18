package db

import (
	"fmt"
	"sync"

	"github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data/v1"
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
	dbFactory data.DataFactory
	once      sync.Once
)

// 这个方法会返回 gorm 连接
// 这个方法应该返回的是全局的一个变量，如果一开始的时候没有初始化好，那么就初始化一次，
// 后续直接拿到这个变量就可以了
func GetDBFactoryOr(mysqlOpts *options.MysqlOptions) (data.DataFactory, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}

	var err error
	once.Do(func() {
		logger := gormx.New( gormx.Config{
			LogLevel: logger.LogLevel(mysqlOpts.LogLevel),
		})

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlOpts.Username, mysqlOpts.Password, mysqlOpts.Host, mysqlOpts.Port, mysqlOpts.Database)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

		dbFactory = &mysqlFactory{
			db: db,
		}
	})
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory, err: %w", err)
	}
	return dbFactory, nil
}

var _ data.DataFactory = (*mysqlFactory)(nil)

type mysqlFactory struct {
	db *gorm.DB
}

// Begin implements data.DataFactory.
func (m *mysqlFactory) Begin() *gorm.DB {
	return m.db.Begin()
}

// Banner implements data.DataFactory.
func (m *mysqlFactory) Banner() data.BannerStore {
	return newBanners(m)
}

// Brands implements data.DataFactory.
func (m *mysqlFactory) Brands() data.BrandsStore {
	return newBrands(m)
}

// CategoryBrands implements data.DataFactory.
func (m *mysqlFactory) CategoryBrands() data.GoodsCategoryBrandStore {
	return newCategoryBrands(m)
}

// Categorys implements data.DataFactory.
func (m *mysqlFactory) Categorys() data.CategoryStore {
	return newCategory(m)
}

// Goods implements data.DataFactory.
func (m *mysqlFactory) Goods() data.GoodsStore {
	return newGoods(m)
}
