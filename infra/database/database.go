package database

import (
	"context"

	"github.com/upvotes-grpc/garcia-paulo/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Ctx      context.Context
	Database *mongo.Database
	Config   *config.Config
}

func NewDatabase(config *config.Config) *Database {
	database := &Database{}
	clientOptions := options.Client().ApplyURI(config.DBSource)
	client, err := mongo.Connect(database.Ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(database.Ctx, nil)
	if err != nil {
		panic(err)
	}

	database.Database = client.Database("upvoteDB")
	return database
}
