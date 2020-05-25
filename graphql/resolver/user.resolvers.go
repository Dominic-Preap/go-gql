package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/my/app/model"
	"github.com/my/app/service"
)

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.Service.User.FindAll(&service.UserFindAll{})
}
