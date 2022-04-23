package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
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
	err := u.users.FindOne(u.ctx, bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (u *UserRepository) CreateUser(user *models.User) error {
	if _, err := u.FindUserByUsername(user.Username); err == nil {
		return status.Errorf(codes.AlreadyExists, "user with username %s already exists", user.Username)
	}
	_, err := u.users.InsertOne(u.ctx, user)
	return err
}
