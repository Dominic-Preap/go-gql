package service

import (
	"github.com/jinzhu/gorm"
	"github.com/my/app/model"
)

// UserService **
type UserService struct {
	DB *gorm.DB
}

// UserFilter --
type UserFilter struct {
	IDs    []int
	Name   *string
	Email  *string
	AgeGte *int

	Limit  *int
	Offset *int
}

// FindAll xxx
func (s *UserService) FindAll(f *UserFilter) ([]*model.User, error) {
	// log.Printf("user opt: %+v", f)

	users := []*model.User{}

	t := s.filter(f)
	t.Order("id ASC").Find(&users)

	return users, nil
}

// FindOne xxx
func (s *UserService) FindOne(f *UserFilter) (*model.User, error) {

	user := &model.User{}
	t := s.filter(f).First(&user)
	return user, GormError(t)
}

func (s *UserService) filter(f *UserFilter) *gorm.DB {
	return s.DB.Scopes(
		WhereSliceInt("id IN (?)", f.IDs),
		WhereString("name = ?", f.Name),
		WhereString("email = ?", f.Email),
		WhereInt("age >= ?", f.AgeGte),
		Limit(f.Limit),
		Offset(f.Offset),
	)
}
