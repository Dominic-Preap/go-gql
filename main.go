package main

import (
	"github.com/my/app/foo"
	"github.com/my/app/service"
)

func main() {
	env := foo.LoadEnv()
	db := foo.ConnectDB(env)
	svc := service.InitService(db)
	foo.InitServer(env, svc)
}

// func main() {
// 	env := foo.LoadEnv()
// 	db := foo.ConnectDB(env)
// 	service := service.InitService(db)

// 	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
// 		Service: service,
// 	}}))

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", dataloader.Middleware(service, srv))

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", env.Port)
// 	log.Fatal(http.ListenAndServe(":"+env.Port, nil))
// }
