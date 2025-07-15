package gormx

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID int64 `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt
}
