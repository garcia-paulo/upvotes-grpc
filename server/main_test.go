package main

import (
	"context"
	"testing"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/database"
	"github.com/garcia-paulo/upvotes-grpc/server/presentation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
}

func TestCreateUser(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()
	users.DeleteOne(ctx, bson.M{"username": "test.user"})

	t.Run("CreateUserSuccess", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "test.user",
			Password: "testPassword",
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("CreateUserDuplicateFail", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, &gen.UserRequest{
			Username: "test.user",
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
			Username: "test.user",
			Password: "",
		})

		if err == nil {
			t.Errorf("Error: invalid password should not be created")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
}

func TestLogin(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.user",
		Password: "testPassword",
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, &gen.UserRequest{
			Username: "test.user",
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
			Username: "test.user",
			Password: "",
		})

		if err == nil {
			t.Errorf("Error: invalid password should not be accepted")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
}

func TestGetPosts(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()

	t.Run("GetPostsSuccess", func(t *testing.T) {
		_, err := server.PostServer.GetPosts(ctx, &gen.Empty{})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})
}

func TestCreatePost(t *testing.T) {
	InitTestSetup()
	ctx := context.WithValue(context.Background(), "username", "test.user")
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.user",
		Password: "testPassword",
	})

	t.Run("CreatePostSuccess", func(t *testing.T) {
		_, err := server.PostServer.CreatePost(ctx, &gen.PostRequest{
			Title: "testTitle",
			Body:  "testContent",
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("CreatePostInvalidTitleFail", func(t *testing.T) {
		_, err := server.PostServer.CreatePost(ctx, &gen.PostRequest{
			Title: "",
			Body:  "testContent",
		})

		if err == nil {
			t.Errorf("Error: invalid title should not be accepted")
		}
	})

	t.Run("CreatePostInvalidBodyFail", func(t *testing.T) {
		_, err := server.PostServer.CreatePost(ctx, &gen.PostRequest{
			Title: "testTitle",
			Body:  "",
		})

		if err == nil {
			t.Errorf("Error: invalid content should not be accepted")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
}

func TestToggleUpvote(t *testing.T) {
	InitTestSetup()
	ctx := context.WithValue(context.Background(), "username", "test.user")
	ctx2 := context.WithValue(context.Background(), "username", "test.user2")
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	users.DeleteOne(ctx, bson.M{"username": "test.user2"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.user",
		Password: "testPassword",
	})
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.user2",
		Password: "testPassword",
	})
	post, err := server.PostServer.CreatePost(ctx, &gen.PostRequest{
		Title: "testTitle",
		Body:  "testContent",
	})

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Run("ToggleUpvoteSuccess", func(t *testing.T) {
		post, err = server.PostServer.ToggleUpvote(ctx, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(post.Upvotes) != 1 {
			t.Errorf("Error: upvote should be added. excepted 1, got %v", len(post.Upvotes))
		}

		post, err = server.PostServer.ToggleUpvote(ctx2, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(post.Upvotes) != 2 {
			t.Errorf("Error: upvote should be added. excepted 2, got %v", len(post.Upvotes))
		}

		post, err = server.PostServer.ToggleUpvote(ctx2, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(post.Upvotes) != 1 {
			t.Errorf("Error: upvote should be added. excepted 1, got %v", len(post.Upvotes))
		}

		post, err = server.PostServer.ToggleUpvote(ctx, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(post.Upvotes) != 0 {
			t.Errorf("Error: upvote should be added. excepted 0, got %v", len(post.Upvotes))
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	users.DeleteOne(ctx, bson.M{"username": "test.user2"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
}
