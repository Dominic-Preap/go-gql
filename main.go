package main

import (
	"log"

	"github.com/my/app/server"
	"github.com/my/app/server/config"
	"github.com/my/app/server/ioredis"
	"github.com/my/app/service"
)

func main() {
	env := config.LoadEnv()
	client := ioredis.InitRedis(env)
	db := server.ConnectDB(env)
	svc := service.InitService(db)

	s := &config.Server{
		Env:      env,
		Database: db,
		Client:   client,
		Service:  svc,
	}

	go ioredis.InitRedisExpiryPubSub(s)
	r := server.InitServer(s)

	log.Printf("🚀 Server ready at http://localhost:%s/", env.Port)
	r.Run(":" + env.Port)
}
