package app

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/my/app/server/config"
	"github.com/my/app/server/gqlclient"
	"github.com/my/app/server/httpclient"
	"github.com/my/app/service"
)

// Server : All of my handlers (db, redis, etc...), hang off of this Server struct
// so these components can access the configuration data
type Server struct {

	// Configuration using viper and used in most of third party library
	Env *config.EnvConfig

	// GORM database instance in case we want to access to database
	Database *gorm.DB

	// Redis client
	Client *redis.Client

	// Repository functions from database models
	Service *service.Service

	// HTTP and REST client for https://reqres.in
	HTTPClient *httpclient.HTTPClient

	// Graphql client for https://countries.trevorblades.com/
	GQLClient *gqlclient.GQLClient
}
