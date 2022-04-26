package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/server/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository struct {
	users *mongo.Collection
	ctx   context.Context
}

func NewUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{
		users: database.Database.Collection("users"),
		ctx:   database.Ctx,
	}
}

func (u *UserRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := u.users.FindOne(u.ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, status.Errorf(codes.NotFound, "user with username %s not found", username)
	}
	return &user, nil
}

func (u *UserRepository) CreateUser(user *models.User) error {
	if _, err := u.FindUserByUsername(user.Username); err == nil {
		return status.Errorf(codes.AlreadyExists, "user with username %s already exists", user.Username)
	}
	if _, err := u.users.InsertOne(u.ctx, user); err != nil {
		return status.Errorf(codes.Internal, "error creating user: %s", err.Error())
	}
	return nil
}
