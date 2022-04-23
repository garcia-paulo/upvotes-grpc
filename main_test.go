package main

import (
	"context"
	"testing"

	"github.com/garcia-paulo/upvotes-grpc/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
	"github.com/garcia-paulo/upvotes-grpc/presentation"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var (
	server *presentation.Server
	users  *mongo.Collection
	posts  *mongo.Collection
)

func InitTestSetup() {
	server = InitializeServer()
	database := database.NewDatabase(config.NewConfig())
	users = database.Database.Collection("users")
	posts = database.Database.Collection("posts")
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.AuthInterceptor.UnaryInterceptor))
	server.RegisterServers(grpcServer)
}

func TestCreateUser(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()

	t.Run("CreateUserSuccess", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "test.bnVsbA",
			Password: "testPassword",
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("CreateUserDuplicateFail", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "test.bnVsbA",
			Password: "testPassword",
		})

		if err == nil {
			t.Errorf("Error: duclicate user should not be created")
		}
	})

	t.Run("CreateUserInvalidUsernameFail", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "",
			Password: "testPassword",
		})

		if err == nil {
			t.Errorf("Error: invalid username should not be created")
		}
	})

	t.Run("CreateUserInvalidPasswordFail", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "test.bnVsbA",
			Password: "",
		})

		if err == nil {
			t.Errorf("Error: invalid password should not be created")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.bnVsbA"})
}

func TestLogin(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.bnVsbA",
		Password: "testPassword",
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, &gen.UserRequest{
			Username: "test.bnVsbA",
			Password: "testPassword",
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("LoginInvalidUsernameFail", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, &gen.UserRequest{
			Username: "",
			Password: "testPassword",
		})

		if err == nil {
			t.Errorf("Error: invalid username should not be accepted")
		}
	})

	t.Run("LoginInvalidPasswordFail", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, &gen.UserRequest{
			Username: "test.bnVsbA",
			Password: "",
		})

		if err == nil {
			t.Errorf("Error: invalid password should not be accepted")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.bnVsbA"})
}
