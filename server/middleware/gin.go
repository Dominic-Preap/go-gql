package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

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

// UseGinContext create a gin middleware to add its context to the context.Context
func UseGinContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetGinContext function to recover the gin.Context from the context.Context struct
func GetGinContext(ctx context.Context) (*gin.Context, error) {
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
