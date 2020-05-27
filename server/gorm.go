package server

import (
	"log"

	"github.com/jinzhu/gorm"

	// need to import postgres driver here
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/my/app/model"
	"github.com/my/app/server/config"
)

// ConnectDB Connecting to a Database
func ConnectDB(env config.EnvConfig) *gorm.DB {
	db, err := gorm.Open(env.GormDialect, env.GormConnectionDSN)
	if err != nil {
		log.Panicln("[ORM] err: ", err)
	}

	db.LogMode(env.GormLogmode) // Log every SQL command

	if env.GormAutomigrate {
		db.AutoMigrate(&model.User{}, &model.Todo{}) // Automigrate tables
	}

	// Create `uuid-ossp` if using postgres
	// https://github.com/jinzhu/gorm/issues/1887#issuecomment-408228087
	if db.Dialect().GetName() == "postgres" {
		db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	}

	log.Print("[ORM] Database connection initialized.")
	return db
}
