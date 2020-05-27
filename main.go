package main

import (
	"github.com/my/app/server"
	"github.com/my/app/service"
)

func main() {
	env := server.LoadEnv()
	db := server.ConnectDB(env)
	svc := service.InitService(db)
	server.InitServer(env, svc)
}
