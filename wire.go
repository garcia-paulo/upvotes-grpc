//go:build wireinject
// +build wireinject

package main

import (
	"github.com/garcia-paulo/upvotes-grpc/application/interceptors"
	"github.com/garcia-paulo/upvotes-grpc/application/servicers"
	"github.com/garcia-paulo/upvotes-grpc/application/token"
	"github.com/garcia-paulo/upvotes-grpc/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/presentation"
	"github.com/garcia-paulo/upvotes-grpc/presentation/servers"
	"github.com/google/wire"
)

func InitializeServer() *presentation.Server {
	panic(wire.Build(
		config.NewConfig,
		database.NewDatabase,
		repositories.NewUserRepository,
		repositories.NewPostRepository,
		servicers.NewUserServicer,
		servicers.NewPostServicer,
		token.NewTokenMaker,
		interceptors.NewAuthInterceptor,
		servers.NewUserServer,
		servers.NewPostServer,
		presentation.NewServer,
	))
}
