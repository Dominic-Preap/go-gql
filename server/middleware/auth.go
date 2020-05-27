package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/my/app/server/config"
	"github.com/my/app/service"

	jwt "github.com/appleboy/gin-jwt/v2"
)

/*
|--------------------------------------------------------------------------
| Auth JWT Middleware
|--------------------------------------------------------------------------
|
| This is a middleware for Gin framework. It uses jwt-go to provide a jwt
| authentication middleware. https://github.com/appleboy/gin-jwt
|
*/

type key string

// login body
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User auth user struct
type User struct {
	ID        string
	UserName  string
	FirstName string
	LastName  string
	Role      string
}

// CurrentUserKey Auth user context key
const CurrentUserKey key = "CurrentUserKey"

// UseAuthJWT JWT middleware for gin
func UseAuthJWT(env config.EnvConfig, svc *service.Service) *jwt.GinJWTMiddleware {

	m, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(env.SecretKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,

		// TokenLookup is a string in the form of "<source>:<name>" that is used to extract token from the request.
		TokenLookup: "header: Authorization",

		// TokenHeadName is a string in the header
		TokenHeadName: "Bearer",

		// ====================================================
		// LOGIN HANDLER SECTION
		// ====================================================

		// Handle login logic
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var body login

			if err := c.ShouldBind(&body); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// ! CHECK DB HERE
			userName := body.Username
			password := body.Password

			if (userName == "admin" && password == "admin") || (userName == "test" && password == "test") {
				return &User{
					ID:        strconv.Itoa(1),
					UserName:  userName,
					FirstName: "Dominic",
					LastName:  "Preap",
					Role:      "ADMIN",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		// Add addition info into json token, in this case, the User response from Authenticator above
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*User); ok {
				mc := jwt.MapClaims{}
				u, _ := json.Marshal(user) // connvert struct to byte
				json.Unmarshal(u, &mc)     // then convert back to map
				return mc
			}
			return jwt.MapClaims{}
		},

		// ====================================================
		// AUTHENTICATION SECTION
		// ====================================================

		// Decode token and then parse into User struct
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			user := &User{}
			mapstructure.Decode(claims, user) // convert map to struct
			return user
		},

		// Think of it as guard logic, might check with redis or database here
		Authorizator: func(data interface{}, c *gin.Context) bool {

			ctx := context.WithValue(c.Request.Context(), CurrentUserKey, data.(*User))
			c.Request = c.Request.WithContext(ctx)

			// ! MIGHT CHECK WITH DB HERE
			if v, ok := data.(*User); ok && v.UserName == "admin" {

				return true
			}
			return false
		},

		// Resolve unauthorize response
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return m
}

// AuthUserContext function to recover auth user from context. used mostly at Resolver level
func AuthUserContext(ctx context.Context) (*User, error) {
	err := errors.New("No user in context")

	if ctx.Value(CurrentUserKey) == nil {
		return nil, err
	}

	user, ok := ctx.Value(CurrentUserKey).(*User)
	if !ok {
		return nil, err
	}

	return user, nil
}
