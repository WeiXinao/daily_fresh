package category

import (
	"database/sql"

	"github.com/WeiXinao/daily_fresh/app/pkg/gormx"
	"github.com/WeiXinao/xkit/slice"
	"gorm.io/gorm"
)

type CategoryDao interface {
}

var _ CategoryDao = (*gormCategoryDao)(nil)

type gormCategoryDao struct {
	db *gorm.DB
}

func NewGormCategoryDao(db *gorm.DB) CategoryDao {
	db.AutoMigrate(&CategoryModel{})
	return &gormCategoryDao{
		db: db,
	}
}

func (gcd *gormCategoryDao)  Transaction(fc func(txDao *gormCategoryDao) error, opts ...*gormx.TxOptions) {
	sqlOpts := []*sql.TxOptions{}
	if len(opts) > 0 {
		sqlOpts = slice.Map(opts, func(idx int, src *gormx.TxOptions) *sql.TxOptions {
			return &sql.TxOptions{
				Isolation: sql.IsolationLevel(src.Isolation),
				ReadOnly: src.ReadOnly,
			}
		})
	}
	gcd.db.Transaction(func(tx *gorm.DB) error {
		dao := &gormCategoryDao{
			db: tx,
		}
		return fc(dao)
	}, sqlOpts...)	
}