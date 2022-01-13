package db

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt time.Time      `json:"created_at" format:"date-time"`
	UpdatedAt time.Time      `json:"updated_at" format:"date-time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true" `
}

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)
