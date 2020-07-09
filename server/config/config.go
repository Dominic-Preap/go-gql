package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	// ProductionEnv indicates env mode is production
	ProductionEnv = "production"

	// DevelopmentEnv indicates env mode is development
	DevelopmentEnv = "development"
)

// BindEnv .
var BindEnv = []string{
	"ENVIRONMENT",
	"PORT",
	"SECRET_KEY",
	"GORM_AUTOMIGRATE",
	"GORM_LOGMODE",
	"GORM_DIALECT",
	"GORM_CONNECTION_DSN",
	"REDIS_ADDRESS",
	"REDIS_PASSWORD",
	"MQTT_HOST",
	"MQTT_USER",
	"MQTT_PASS",
}

// EnvConfig all configuration for the server are define here
type EnvConfig struct {
	// General Config
	// --------------
	Environment string `validate:"required,oneof=development production" mapstructure:"ENVIRONMENT"`
	Port        string `validate:"required" mapstructure:"PORT"`
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

	// Redis Config
	// ---------------
	MQTTHost string `validate:"required" mapstructure:"MQTT_HOST"`
	MQTTUser string `validate:"required" mapstructure:"MQTT_USER"`
	MQTTPass string `validate:"required" mapstructure:"MQTT_PASS"`
}

// LoadEnv Load environment variable from .env file
func LoadEnv() (*EnvConfig, error) {
	c := &EnvConfig{}

	v := viper.New() // Create a new viper instance
	if os.Getenv("ENVIRONMENT") == DevelopmentEnv {
		v.SetConfigFile(".env") // name of config file (without extension)
		v.SetConfigType("env")  // if the config file does not have the extension in the name
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("Error reading config file, %s", err)
		}
	} else {
		v.AutomaticEnv()
		for _, key := range BindEnv {
			v.BindEnv(key)
		}
	}

	if err := v.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("Unable to decode config into struct, %s", err)
	}

	if err := validator.New().Struct(c); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf(`Error:Field "%[2]v" on "%[1]v"`, e.StructField(), e.ActualTag())
		}
		return nil, fmt.Errorf("Please check your .env file again")
	}

	return c, nil
}
