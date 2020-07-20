package service

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"github.com/my/app/model"
)

// UserService .
type UserService struct {
	DB *gorm.DB
}

// UserFilter .
type UserFilter struct {
	IDs    []int
	Name   *string
	Email  *string
	AgeGte *int

	Limit  *int
	Offset *int
}

// FindAll .
func (s *UserService) FindAll(f *UserFilter) ([]*model.User, error) {
	// log.Printf("user opt: %+v", f)
	q := sq.Select("*").From("users")
	q = s.filter(q, f)
	q = LimitOffset(q, f.Limit, f.Offset)

	users := []*model.User{}
	t := s.DB.Raw(q.MustSql()).Find(&users)
	return users, GormError(t)
}

// FindOne .
func (s *UserService) FindOne(f *UserFilter) (*model.User, error) {
	q := sq.Select("*").From("users").Limit(1)
	q = s.filter(q, f)

	user := &model.User{}
	t := s.DB.Raw(q.MustSql()).First(&user)
	return user, GormError(t)
}

func (s *UserService) filter(q sq.SelectBuilder, f *UserFilter) sq.SelectBuilder {
	return q.PlaceholderFormat(sq.Dollar).
		Where(Cond(sq.Eq{"id": f.IDs}, IsNull(f.IDs))).
		Where(Cond(sq.Eq{"name": f.Name}, IsNullOrEmpty(f.Name))).
		Where(Cond(sq.Eq{"email": f.Email}, IsNullOrEmpty(f.Email))).
		Where(Cond(sq.GtOrEq{"age": f.AgeGte}, IsNull(f.AgeGte))).
		OrderBy("id ASC")
}
