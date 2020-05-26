package foo

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/my/app/graphql/dataloader"
	"github.com/my/app/graphql/generated"
	"github.com/my/app/graphql/resolver"
	"github.com/my/app/service"
)

type key string

const (
	// GinContextKey context key for Gin Http Server
	GinContextKey key = "GinContextKey"
)

// InitServer .
func InitServer(env EnvConfig, svc *service.Service) {
	gin.SetMode(gin.ReleaseMode)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()
	r.Use(ginContextToContextMiddleware())
	r.POST("/query", graphqlHandler(svc))
	r.GET("/", playgroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", env.Port)
	r.Run(":" + env.Port)
}

// GinContextFromContext Define a function to recover the `gin.Context` from the `context.Context` struct
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

// GinContextToContextMiddleware .
func ginContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Defining the Graphql handler
func graphqlHandler(svc *service.Service) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := generated.Config{Resolvers: &resolver.Resolver{Service: svc}}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		// h.ServeHTTP(c.Writer, c.Request)
		dataloader.Middleware(svc, h).ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
