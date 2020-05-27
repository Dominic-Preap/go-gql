package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/my/app/graphql/dataloader"
	"github.com/my/app/graphql/generated"
	modelgen "github.com/my/app/graphql/model"
	"github.com/my/app/graphql/resolver"
	"github.com/my/app/service"
)

type key string

// InitServer start gin server
func InitServer(env EnvConfig, svc *service.Service) {
	if env.Env == ProductionEnv {
		gin.SetMode(gin.ReleaseMode) // set gin mode to release mode on production env
	}

	r := gin.Default()            // create gin router with logger and recovery middleware
	r.Use(corsMiddleware())       // use CORS middleware
	r.Use(ginContextMiddleware()) // use gin context middleware for graphql context

	auth := authMiddleware(env, svc) // init auth middleware that contain login handler and refresh token
	r.GET("/refresh_token", auth.RefreshHandler)
	r.POST("/login", auth.LoginHandler)
	r.POST("/query", auth.MiddlewareFunc(), graphqlHandler(svc))

	if env.Env != ProductionEnv {
		r.GET("/", playgroundHandler()) // Graphql Playground does not avaliable on production
	}

	log.Printf("🚀 Server ready at http://localhost:%s/", env.Port)
	r.Run(":" + env.Port)
}

// Defining the Graphql handler
func graphqlHandler(svc *service.Service) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := generated.Config{Resolvers: &resolver.Resolver{Service: svc}}
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

// Defining CORS middleware
func corsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	return cors.New(config)
}

/*
|--------------------------------------------------------------------------
| Accessing gin.Context
|--------------------------------------------------------------------------
|
| At the Resolver level, gqlgen gives you access to the context.Context object.
| One way to access the gin.Context is to add it to the context and retrieve it again.
|
| https://gqlgen.com/recipes/gin/#accessing-gincontext
*/

// GinContextKey context key for Gin Http Server
const GinContextKey key = "GinContextKey"

// create a gin middleware to add its context to the context.Context
func ginContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext function to recover the gin.Context from the context.Context struct
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
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
	if u, err := AuthUserContext(ctx); err != nil || u.Role != role.String() {
		return nil, errors.New("Access denied")
	}
	return next(ctx)
}
