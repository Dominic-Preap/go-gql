<img src="https://raw.githubusercontent.com/egonelbre/gophers/ac77b513f41f44a7805694063aaef16ccd95a9b3/vector/party/birthday.svg" alt="logo" align="right" width="140" />

# Graphql Go Boilerplate

## 📋 Introduction

A simple Graphql Go HTTP server boilerplate. This boilerplate project can be fully customized to your needs.

Thanks to VSCode and Go, we can have a fully intellisense when coding. It come with `.vscode/extensions.json` for recommended extensions.

## Requirements

1. Download and install [Go](https://golang.org/dl/)
2. Run the command `go get -u` to install dependencies
3. (Opt) Install [Air](https://github.com/cosmtrek/air) if you want ☁️ Live reload

```sh
# install air via go commannd
$ go get -u github.com/cosmtrek/air
```

4. Copy `.env.example` to new env file `.env` and fill the environment variables

```sh
# Application (development, production)
ENV=development
PORT=8080
SECRET_KEY="5cR3t_K37"

# Gorm Credential
GORM_LOGMODE=true
GORM_DIALECT=postgres
GORM_CONNECTION_DSN="postgres://username:password@localhost:port/database"
```

## Quick start

```sh
# run server and visit localhost:8080 for graphql playground
$ go run main.go

# or run live reload go via air and visit localhost:8080
$ air -d
```

## Features

- [] TODO

## Development

- [] TODO

```sh
# generate dataloader file if you are using one
$ dataloaden UserLoader int *github.com/my/app/model.User

# -or-
$ go run github.com/vektah/dataloaden UserLoader int *github.com/my/app/model.User

# run command everything you're changing the Graphql Schema file *.gql
$ gqlgen generate
$ go run github.com/99designs/gqlgen

# clean your project dependencies
$ go get tidy
```
