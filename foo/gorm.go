package foo

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/my/app/model"
)

// ConnectDB xx
func ConnectDB(env EnvConfig) *gorm.DB {
	db, err := gorm.Open(env.GormDialect, env.GormConnectionDSN)
	if err != nil {
		log.Panicln("[ORM] err: ", err)
	}

	// Log every SQL command on dev, @prod: this should be disabled? Maybe.
	db.LogMode(env.GormLogmode)

	// Automigrate tables
	if env.GormAutomigrate {
		db.AutoMigrate(&model.User{}, &model.Todo{})
	}

	log.Print("[ORM] Database connection initialized.")
	return db
}
