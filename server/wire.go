//go:build wireinject
// +build wireinject

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
		interceptors.NewUnaryInterceptor,
		servers.NewUserServer,
		servers.NewPostServer,
		presentation.NewServer,
	))
}
