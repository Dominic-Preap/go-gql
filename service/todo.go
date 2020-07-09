package service

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/my/app/model"

	sq "github.com/Masterminds/squirrel"
)

// TodoService .
type TodoService struct {
	DB *gorm.DB
}

// TodoFilter --
type TodoFilter struct {
	ID       *int
	Text     *string
	TextLike *string
	Done     *bool
	UserID   *int
	UserIDs  []int

	Limit  *int
	Offset *int
}

// FindAll .
func (s *TodoService) FindAll(f *TodoFilter) ([]*model.Todo, error) {
	q := sq.Select("*").From("todos")
	q = s.filter(q, f)
	q = LimitOffset(q, f.Limit, f.Offset)

	todos := []*model.Todo{}
	t := s.DB.Raw(q.MustSql()).Find(&todos)
	return todos, GormError(t)
}

// Count .
func (s *TodoService) Count(f *TodoFilter) (int, error) {
	q := sq.Select("COUNT(*) AS total").From("todos")
	q = s.filter(q, f)

	x := Total{}
	t := s.DB.Raw(q.MustSql()).Find(&x)
	return x.total, GormError(t)
}

type FindAllXResult struct {
	Data []*model.Todo
	Err  error
}

// FindAllX .
func (s *TodoService) FindAllX(f *TodoFilter) <-chan FindAllXResult {
	r := make(chan FindAllXResult)
	go func() {
		// defer close(r)
		data, err := s.FindAll(f)
		r <- FindAllXResult{data, err}
	}()
	return r

}

// FindAndCountAll .
func (s *TodoService) FindAndCountAll(f *TodoFilter) ([]*model.Todo, int, error) {

	c, _ := s.FindAll(f)
	a, _ := s.FindAll(f)
	x, _ := s.FindAll(f)
	y, _ := s.FindAll(f)
	d, _ := s.FindAll(f)
	e, _ := s.FindAll(f)
	z, _ := s.FindAll(f)

	log.Print(a, c, x, y, z, d, e)

	return []*model.Todo{}, 0, nil
}

// FindOne .
func (s *TodoService) FindOne(f *TodoFilter) (*model.Todo, error) {
	q := sq.Select("*").From("todos").Limit(1)
	q = s.filter(q, f)

	todo := &model.Todo{}
	t := s.DB.Raw(q.MustSql()).First(&todo)
	return todo, GormError(t)
}

// Create .
func (s *TodoService) Create(todo *model.Todo) (*model.Todo, error) {
	err := s.DB.Create(&todo).Error
	return todo, err
}

func (s *TodoService) filter(q sq.SelectBuilder, f *TodoFilter) sq.SelectBuilder {
	return q.PlaceholderFormat(sq.Dollar).
		Where(sq.NotEq{"id": nil}).
		Where(Cond(sq.Eq{"id": f.ID}, IsNull(f.ID))).
		Where(Cond(sq.Eq{"done": f.Done}, IsNull(f.Done))).
		Where(Cond(sq.Eq{"text": f.Text}, IsNullOrEmpty(f.Text))).
		Where(Cond(sq.ILike{"text": Like(f.TextLike)}, IsNullOrEmpty(f.TextLike))).
		Where(Cond(sq.Eq{"user_id": f.UserID}, IsNull(f.UserID))).
		Where(Cond(sq.Eq{"user_id": f.UserIDs}, IsNull(f.UserIDs))).
		OrderBy("id ASC")
}
