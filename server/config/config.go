package config

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/my/app/service"
	"github.com/spf13/viper"
)

const (
	// ProductionEnv indicates env mode is production
	ProductionEnv = "production"

	// DevelopmentEnv indicates env mode is development
	DevelopmentEnv = "development"
)

// Server : All of my handlers (db, redis, etc...), hang off of this Server struct
// so these components can access the configuration data
type Server struct {

	// Configuration using viper and used in most of third party library
	Env *EnvConfig

	// GORM database instance in case we want to access to database
	Database *gorm.DB

	// Redis client
	Client *redis.Client

	// Repository functions from database models
	Service *service.Service
}

// EnvConfig all configuration for the server are define here
type EnvConfig struct {
	// General Config
	// --------------
	Environment string `validate:"required,oneof=development production"`
	Port        string `validate:"required"`
	SecretKey   string `validate:"required" mapstructure:"SECRET_KEY"`

	// Database Config
	// ---------------
	GormAutomigrate   bool   `mapstructure:"GORM_AUTOMIGRATE"`
	GormLogmode       bool   `mapstructure:"GORM_LOGMODE"`
	GormDialect       string `validate:"required" mapstructure:"GORM_DIALECT"`
	GormConnectionDSN string `validate:"required" mapstructure:"GORM_CONNECTION_DSN"`

	// Redis Config
	// ---------------
	RedisAddress  string `validate:"required" mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `validate:"required" mapstructure:"REDIS_PASSWORD"`
}

// LoadEnv Load environment variable from .env file
func LoadEnv() *EnvConfig {
	c := &EnvConfig{}

	v := viper.New()        // Create a new viper instance
	v.SetConfigFile(".env") // name of config file (without extension)
	v.SetConfigType("env")  // if the config file does not have the extension in the name

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := v.Unmarshal(c); err != nil {
		log.Fatalf("Unable to decode config into struct, %s \n", err)
	}

	if err := validator.New().Struct(c); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf(`Error:Field "%[2]v" on "%[1]v"`, e.StructField(), e.ActualTag())
		}
		log.Fatal("Please check your .env file again")
	}

	return c
}
