// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/garcia-paulo/upvotes-grpc/server/application/interceptors"
	"github.com/garcia-paulo/upvotes-grpc/server/application/servicers"
	"github.com/garcia-paulo/upvotes-grpc/server/application/token"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/database"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/server/presentation"
	"github.com/garcia-paulo/upvotes-grpc/server/presentation/servers"
)

// Injectors from wire.go:

func InitializeServer() *presentation.Server {
	configConfig := config.NewConfig()
	databaseDatabase := database.NewDatabase(configConfig)
	userRepository := repositories.NewUserRepository(databaseDatabase)
	tokenMaker := token.NewTokenMaker(configConfig)
	userServicer := servicers.NewUserServicer(userRepository, tokenMaker)
	userServer := servers.NewUserServer(userServicer)
	postRepository := repositories.NewPostRepository(databaseDatabase, userRepository)
	postServicer := servicers.NewPostServicer(postRepository)
	postServer := servers.NewPostServer(postServicer)
	unaryInterceptor := interceptors.NewUnaryInterceptor(tokenMaker)
	server := presentation.NewServer(userServer, configConfig, postServer, unaryInterceptor)
	return server
}
