package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/my/app/graphql/dataloader"
	"github.com/my/app/graphql/generated"
	modelgen "github.com/my/app/graphql/model"
	"github.com/my/app/model"
	"github.com/my/app/service"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input modelgen.InputTodo) (*model.Todo, error) {
	t := &model.Todo{
		UserID: input.UserID,
		Text:   input.Text,
	}
	return r.Service.Todo.Create(t)
}

func (r *queryResolver) Todos(ctx context.Context, filter *modelgen.TodoFilter, limit *int, offset *int) ([]*model.Todo, error) {
	return r.Service.Todo.FindAll(&service.TodoFilter{
		// UserID: filter.UserID,
		// Done:   filter.Done,
		// TextLike: filter.Text,
	})
}

func (r *queryResolver) Todo(ctx context.Context, id int) (*model.Todo, error) {
	return r.Service.Todo.FindOne(&service.TodoFilter{ID: id})
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return dataloader.For(ctx).UserByID.Load(obj.UserID)
}

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type todoResolver struct{ *Resolver }
