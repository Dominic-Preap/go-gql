package service

import (
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
)

// Total .
type Total struct {
	total int
}

// Service .
type Service struct {
	User UserService
	Todo TodoService
}

// Init .
func Init(db *gorm.DB) *Service {
	return &Service{
		User: UserService{DB: db},
		Todo: TodoService{DB: db},
	}
}

// Cond verify if condition is valid to add to SQL statement
func Cond(pred sq.Sqlizer, byPass bool) sq.Sqlizer {
	if byPass {
		return sq.Expr("")
	}
	return pred
}

// Like .
func Like(arg *string) *string {
	like := ""
	if arg != nil {
		like = "%" + *arg + "%"
	}
	return &like
}

// LimitOffset set limit and offset to SQL if value exist
func LimitOffset(q sq.SelectBuilder, limit *int, offset *int) sq.SelectBuilder {
	if limit != nil {
		q.Limit(uint64(*limit))
	}
	if offset != nil {
		q.Offset(uint64(*offset))
	}
	return q
}

// GormError .
func GormError(t *gorm.DB) error {
	if len(t.GetErrors()) > 0 {
		return t.GetErrors()[0]
	}
	return nil
}

// ==============================
// HELPER
// ==============================

// IsNull .
func IsNull(arg interface{}) bool {
	r := reflect.ValueOf(arg)
	return r.IsNil()
}

// IsNullOrEmpty .
func IsNullOrEmpty(arg *string) bool {
	return arg == nil || *arg == ""
}

// IsEmpty .
func IsEmpty(arg *string) bool {
	return *arg == ""
}
