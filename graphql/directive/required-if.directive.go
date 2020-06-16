package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// RequiredIfDirective .
func RequiredIfDirective(ctx context.Context, obj interface{}, next graphql.Resolver, field string, value interface{}) (res interface{}, err error) {
	v := obj.(map[string]interface{})[field]

	if value == nil && v == nil {
		return nil, fmt.Errorf("There's a field that required %s", field)
	}

	if value != nil && fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
		return nil, fmt.Errorf("Field: %[1]v is required to match value: %[2]v", field, value)
	}

	return next(ctx)
}
