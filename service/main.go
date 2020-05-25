package service

import (
	"github.com/jinzhu/gorm"
)

// Service --
type Service struct {
	User UserService
	Todo TodoService
}

// InitService --
func InitService(db *gorm.DB) *Service {

	return &Service{
		User: UserService{DB: db},
		Todo: TodoService{DB: db},
	}
}

// Where -
func Where(query string, arg interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if arg == "" || arg == 0 || arg == nil {
			return db
		} else {
			return db.Where(query, arg)
		}

		// switch arg.(type) {
		// case bool:
		// 	return cast.ToBool(val)
		// case string:
		// 	return cast.ToString(val)
		// case int32, int16, int8, int:
		// 	return cast.ToInt(val)
		// case uint:
		// 	return cast.ToUint(val)
		// case uint32:
		// 	return cast.ToUint32(val)
		// case uint64:
		// 	return cast.ToUint64(val)
		// case int64:
		// 	return cast.ToInt64(val)
		// case float64, float32:
		// 	return cast.ToFloat64(val)
		// case time.Time:
		// 	return cast.ToTime(val)
		// case time.Duration:
		// 	return cast.ToDuration(val)
		// case []string:
		// 	return cast.ToStringSlice(val)
		// case []int:
		// 	return cast.ToIntSlice(val)
		// }

	}
}
