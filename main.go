package main

import (
	"log"

	"github.com/my/app/server"
	"github.com/my/app/server/app"
	"github.com/my/app/server/config"
	"github.com/my/app/server/gqlclient"
	"github.com/my/app/server/httpclient"
	"github.com/my/app/server/ioredis"
	"github.com/my/app/server/msgbroker"
	"github.com/my/app/service"
)

func main() {
	env := config.LoadEnv()
	db := server.ConnectDB(env)
	svc := service.Init(db)
	api := httpclient.Init(env)
	gql := gqlclient.Init(env)
	client := ioredis.Init(env)
	broker := msgbroker.Init(env)

	s := &app.Server{
		Env:        env,
		Database:   db,
		Client:     client,
		Service:    svc,
		HTTPClient: api,
		GQLClient:  gql,
		MSGBroker:  broker,
	}

	go ioredis.InitRedisExpiryPubSub(s)
	r := server.InitServer(s)

	log.Printf("🚀 Server ready at http://localhost:%s/", env.Port)
	r.Run(":" + env.Port)
}
