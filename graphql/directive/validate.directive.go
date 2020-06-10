package directive

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

// ValidateDirective .
func ValidateDirective(ctx context.Context, obj interface{}, next graphql.Resolver, field string, rules string) (res interface{}, err error) {
	if field == "" || rules == "" {
		return nil, errors.New("Missing validate field or rules")
	}

	value := obj.(map[string]interface{})[field]
	if err := validator.New().Var(value, rules); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf(`Validation failed on field "%[1]v", condition "%[2]v"`, field, err.ActualTag())
		}
	}
	return next(ctx)
}
