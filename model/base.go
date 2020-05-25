package model

import (
	"time"
)

// BaseModel xxx
type BaseModel struct {
	CreatedBy *uint
	UpdatedBy *uint
	CreatedAt *time.Time `gorm:"index;not null;default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"index"`
}
