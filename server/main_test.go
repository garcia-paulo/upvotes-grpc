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

var (
	mockUser = &gen.UserRequest{
		Username: "test.user",
		Password: "testPassword",
	}
	mockPost = &gen.PostRequest{
		Title: "testTitle",
		Body:  "testContent",
	}
)

func TestCreateUser(t *testing.T) {
	InitTestSetup()
	ctx := context.Background()
	users.DeleteOne(ctx, bson.M{"username": "test.user"})

	t.Run("CreateUserSuccess", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, mockUser)

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("CreateUserDuplicateFail", func(t *testing.T) {
		_, err := server.UserServer.CreateUser(ctx, mockUser)

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
	server.UserServer.CreateUser(ctx, mockUser)

	t.Run("LoginSuccess", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, mockUser)

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

	t.Run("LoginUserNotFoundFail", func(t *testing.T) {
		_, err := server.UserServer.Login(ctx, &gen.UserRequest{
			Username: "test.user2",
			Password: "",
		})

		if err == nil {
			t.Errorf("Error: invalid user should not be accepted")
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
	server.UserServer.CreateUser(ctx, mockUser)

	t.Run("CreatePostSuccess", func(t *testing.T) {
		_, err := server.PostServer.CreatePost(ctx, mockPost)

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
	server.UserServer.CreateUser(ctx, mockUser)
	server.UserServer.CreateUser(ctx, &gen.UserRequest{
		Username: "test.user2",
		Password: "testPassword",
	})
	post, err := server.PostServer.CreatePost(ctx, mockPost)

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

	t.Run("ToggleUpvoteInvalidPostIdFail", func(t *testing.T) {
		_, err := server.PostServer.ToggleUpvote(ctx, &gen.PostIdRequest{
			PostId: "",
		})

		if err == nil {
			t.Errorf("Error: invalid post id should not be accepted")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	users.DeleteOne(ctx, bson.M{"username": "test.user2"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
}

func TestDeletePost(t *testing.T) {
	InitTestSetup()
	ctx := context.WithValue(context.Background(), "username", "test.user")
	ctx2 := context.WithValue(context.Background(), "username", "test.user2")
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
	server.UserServer.CreateUser(ctx, mockUser)
	post, err := server.PostServer.CreatePost(ctx, mockPost)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Run("DeletePostSuccess", func(t *testing.T) {
		_, err = server.PostServer.DeletePost(ctx, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("DeletePostFail", func(t *testing.T) {
		_, err = server.PostServer.DeletePost(ctx2, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err == nil {
			t.Errorf("Error: delete post should fail")
		}
	})

	t.Run("DeletePostPermissionDeniedFail", func(t *testing.T) {
		_, err = server.PostServer.DeletePost(ctx2, &gen.PostIdRequest{
			PostId: post.Id,
		})

		if err == nil {
			t.Errorf("Error: delete post should fail")
		}
	})

	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
}

func TestUpdatePost(t *testing.T) {
	InitTestSetup()
	ctx := context.WithValue(context.Background(), "username", "test.user")
	ctx2 := context.WithValue(context.Background(), "username", "test.user2")
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
	server.UserServer.CreateUser(ctx, mockUser)
	post, err := server.PostServer.CreatePost(ctx, mockPost)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Run("UpdatePostSuccess", func(t *testing.T) {
		_, err = server.PostServer.UpdatePost(ctx, &gen.PostUpdateRequest{
			PostId: post.Id,
			Title:  "Testes",
			Body:   "Testes",
		})

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("UpdatePostInvalidTitleFail", func(t *testing.T) {
		_, err = server.PostServer.UpdatePost(ctx, &gen.PostUpdateRequest{
			PostId: post.Id,
			Title:  "",
			Body:   "Testes",
		})

		if err == nil {
			t.Errorf("Error: update post should fail")
		}
	})

	t.Run("UpdatePostInvalidBodyFail", func(t *testing.T) {
		_, err = server.PostServer.UpdatePost(ctx, &gen.PostUpdateRequest{
			PostId: post.Id,
			Title:  "Testes",
			Body:   "",
		})

		if err == nil {
			t.Errorf("Error: update post should fail")
		}
	})

	t.Run("UpdatePostPermissionDeniedFail", func(t *testing.T) {
		_, err = server.PostServer.UpdatePost(ctx2, &gen.PostUpdateRequest{
			PostId: post.Id,
			Title:  "Testes",
			Body:   "Testes",
		})

		if err == nil {
			t.Errorf("Error: update post should fail")
		}
	})
	users.DeleteOne(ctx, bson.M{"username": "test.user"})
	posts.DeleteMany(ctx, bson.M{"author": "test.user"})
}
