//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/upvotes-grpc/garcia-paulo/infra/config"
	"github.com/upvotes-grpc/garcia-paulo/infra/database"
	"github.com/upvotes-grpc/garcia-paulo/infra/repositories"
	"github.com/upvotes-grpc/garcia-paulo/presentation"
	"github.com/upvotes-grpc/garcia-paulo/presentation/servers"
)

func InitializeServer() *presentation.Server {
	panic(wire.Build(
		config.NewConfig,
		database.NewDatabase,
		repositories.NewUserRepository,
		servers.NewUserServer,
		presentation.NewServer,
	))
}
