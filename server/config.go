package server

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

const (
	// ProductionEnv indicates env mode is production
	ProductionEnv = "production"

	// DevelopmentEnv indicates env mode is development
	DevelopmentEnv = "development"
)

// EnvConfig all configuration for the server are define here
type EnvConfig struct {
	// General Config
	// --------------
	Env       string `validate:"required,oneof=development production"`
	Port      string `validate:"required"`
	SecretKey string `validate:"required" mapstructure:"SECRET_KEY"`

	// Database Config
	// ---------------
	GormAutomigrate   bool   `mapstructure:"GORM_AUTOMIGRATE"`
	GormSeedDb        bool   `mapstructure:"GORM_SEED_DB"`
	GormLogmode       bool   `mapstructure:"GORM_LOGMODE"`
	GormDialect       string `validate:"required" mapstructure:"GORM_DIALECT"`
	GormConnectionDSN string `validate:"required" mapstructure:"GORM_CONNECTION_DSN"`

	// Firebase Config
	// ---------------
	FirebaseDBUrl          string `validate:"required" mapstructure:"FIREBASE_DB_URL"`
	FirebaseCredentialFile string `validate:"required" mapstructure:"FIREBASE_CREDENTIAL_FILE"`
}

// LoadEnv Load environment variable from .env file
func LoadEnv() EnvConfig {
	c := EnvConfig{}

	v := viper.New()        // Create a new viper instance
	v.SetConfigFile(".env") // name of config file (without extension)
	v.SetConfigType("env")  // if the config file does not have the extension in the name

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := v.Unmarshal(&c); err != nil {
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
