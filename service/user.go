package service

import (
	"github.com/jinzhu/gorm"
	"github.com/my/app/model"
)

// UserService **
type UserService struct {
	DB *gorm.DB
}

// UserFindAll --
type UserFindAll struct {
	IDs    []int
	Name   string
	Email  string
	AgeGte *int
}

// FindAll xxx
func (u *UserService) FindAll(f *UserFindAll) ([]*model.User, error) {
	// log.Printf("user opt: %+v", f)

	users := []*model.User{}

	t := u.DB.Scopes(
		WhereSliceInt("id IN (?)", f.IDs),
		WhereString("name = ?", f.Name),
		WhereString("email = ?", f.Email),
		WhereInt("age >= ?", f.AgeGte),
	)

	t = t.Order("id ASC").Find(&users)

	return users, nil
}

// FindOne xxx
func (u *UserService) FindOne(f *UserFindAll) (*model.User, error) {
	user := &model.User{}
	u.DB.First(&user)
	return user, nil
}
