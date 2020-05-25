package model

import (
	"database/sql"
	"time"
)

// User -
type User struct {
	ID           int    `json:"id" gorm:"primary_key,auto_increment"`
	Name         string `json:"name"`
	Age          sql.NullInt64
	Birthday     *time.Time
	Email        string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Role         string  `gorm:"size:255"`        // set field size to 255
	MemberNumber *string `gorm:"unique;not null"` // set member number to unique and not null
	Num          int     `gorm:"AUTO_INCREMENT"`  // set num to auto incrementable
	Address      string  `gorm:"index:addr"`      // create index with name `addr` for address
	IgnoreMe     int     `gorm:"-"`               // ignore this field

	BaseModel
}
