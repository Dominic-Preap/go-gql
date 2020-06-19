package service

import (
	"log"

	"github.com/jinzhu/gorm"
)

// Service --
type Service struct {
	User UserService
	Todo TodoService
}

// Init --
func Init(db *gorm.DB) *Service {
	return &Service{
		User: UserService{DB: db},
		Todo: TodoService{DB: db},
	}
}

// GormError .
func GormError(t *gorm.DB) error {
	if len(t.GetErrors()) > 0 {
		return t.GetErrors()[0]
	}
	return nil
}

// Limit .
func Limit(arg *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == nil || *arg <= 0 {
			return db
		}
		return db.Limit(*arg)
	}
}

// Offset .
func Offset(arg *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == nil || *arg <= 0 {
			return db
		}
		return db.Offset(*arg)
	}
}

// WhereString .
func WhereString(query string, arg string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == "" || arg == "%%" {
			return db
		}
		return db.Where(query, arg)
	}
}

// WhereInt .
func WhereInt(query string, arg *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == nil {
			return db
		}
		return db.Where(query, arg)
	}
}

// WhereBool .
func WhereBool(query string, arg *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == nil {
			return db
		}
		return db.Where(query, arg)
	}
}

// WhereSliceInt .
func WhereSliceInt(query string, arg []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(arg) <= 0 || arg == nil {
			return db
		}
		return db.Where(query, arg)
	}
}

// WhereSliceString .
func WhereSliceString(query string, arg interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v, err := arg.([]string)
		if err {
			log.Printf("Error Query: %[1]v, %[2]v", query, arg)
		}

		if len(v) <= 0 || arg == nil {
			return db
		}
		return db.Where(query, arg)
	}
}
