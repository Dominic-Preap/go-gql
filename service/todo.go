package service

import (
	"github.com/jinzhu/gorm"
	"github.com/my/app/model"
)

// TodoService .
type TodoService struct {
	DB *gorm.DB
}

// TodoFilter --
type TodoFilter struct {
	ID       *int
	Text     string
	TextLike string
	Done     *bool
	UserID   *int
	UserIDs  []int

	Limit  *int
	Offset *int
}

// FindAll .
func (s *TodoService) FindAll(f *TodoFilter) ([]*model.Todo, error) {
	todos := []*model.Todo{}

	t := s.filter(f)
	t.Order("id ASC").Find(&todos)

	return todos, nil
}

// FindOne .
func (s *TodoService) FindOne(f *TodoFilter) (*model.Todo, error) {
	todo := &model.Todo{}

	t := s.filter(f)
	t.First(&todo)

	return todo, nil
}

// Create .
func (s *TodoService) Create(todo *model.Todo) (*model.Todo, error) {
	err := s.DB.Create(&todo).Error
	return todo, err
}

func (s *TodoService) filter(f *TodoFilter) *gorm.DB {
	return s.DB.Scopes(
		WhereInt("id = ?", f.ID),
		WhereBool("done = ?", f.Done),
		WhereString("text = ?", f.Text),
		WhereString("text LIKE ?", "%"+f.TextLike+"%"),
		WhereInt("user_id = ?", f.UserID),
		WhereSliceInt("user_id IN (?)", f.UserIDs),
		Limit(f.Limit),
		Offset(f.Offset),
	)
}
