package main

import (
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

	go ioredis.InitRedisExpiryPubSub(svc, client)
	server.InitServer(env, svc, client)
}
