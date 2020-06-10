package directive

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	modelgen "github.com/my/app/graphql/model"
	"github.com/my/app/server/middleware"
)

/*
|--------------------------------------------------------------------------
| Custom Directives
|--------------------------------------------------------------------------
|
| A directive is an identifier preceded by a @ character. Used to authorization
| via GraphQL Schema or any somesort of validation directly on Graphql Schema.
|
| https://www.apollographql.com/docs/graphql-tools/schema-directives/
| https://www.apollographql.com/docs/apollo-server/security/authentication
*/

// AuthDirective Authorize GraphQL Schema via '@auth' directive,
// so we don't have to validate everything at the Resolver level
func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role *modelgen.Role) (interface{}, error) {
	if u, err := middleware.GetAuthUser(ctx); err != nil || u.Role != role.String() {
		return nil, errors.New("Access denied")
	}
	return next(ctx)
}
