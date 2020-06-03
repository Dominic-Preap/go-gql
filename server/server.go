package server

import (
	"context"
	"errors"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"github.com/my/app/graphql/dataloader"
	"github.com/my/app/graphql/generated"
	modelgen "github.com/my/app/graphql/model"
	"github.com/my/app/graphql/resolver"
	"github.com/my/app/server/config"
	"github.com/my/app/server/middleware"
	"github.com/my/app/service"
)

// InitServer start gin server
func InitServer(env config.EnvConfig, svc *service.Service, client *redis.Client) {
	if env.Env == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode) // set gin mode to release mode on production env
	}

	r := gin.Default()                // create gin router with logger and recovery middleware
	r.Use(middleware.UseCors())       // use CORS middleware
	r.Use(middleware.UseGinContext()) // use gin context middleware for graphql context

	auth := middleware.UseAuthJWT(env, svc) // init auth middleware that contain login handler and refresh token
	r.GET("/refresh_token", auth.RefreshHandler)
	r.POST("/login", auth.LoginHandler)
	r.POST("/query", auth.MiddlewareFunc(), graphqlHandler(svc, client))

	if env.Env != config.ProductionEnv {
		r.GET("/", playgroundHandler()) // Graphql Playground does not avaliable on production
	}

	log.Printf("🚀 Server ready at http://localhost:%s/", env.Port)
	r.Run(":" + env.Port)
}

// Defining the Graphql handler
func graphqlHandler(svc *service.Service, client *redis.Client) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := generated.Config{Resolvers: &resolver.Resolver{Service: svc, Client: client}}
	c.Directives.Auth = authDirective // implement auth directive

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		// bind dataloader middleware on Graphql server
		dataloader.Middleware(svc, h).ServeHTTP(c.Writer, c.Request)
		// h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

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

// Authorize GraphQL Schema via '@auth' directive,
// so we don't have to validate everything at the Resolver level
func authDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role *modelgen.Role) (interface{}, error) {
	if u, err := middleware.GetAuthUser(ctx); err != nil || u.Role != role.String() {
		return nil, errors.New("Access denied")
	}
	return next(ctx)
}
