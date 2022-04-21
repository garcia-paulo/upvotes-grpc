package repositories

import (
	"context"

	"github.com/upvotes-grpc/garcia-paulo/domain/models"
	"github.com/upvotes-grpc/garcia-paulo/infra/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Users *mongo.Collection
	Ctx   context.Context
}

func NewUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{
		Users: database.Database.Collection("users"),
		Ctx:   database.Ctx,
	}
}

func (u *UserRepository) CreateUser(user *models.User) error {
	result, err := u.Users.InsertOne(u.Ctx, user)
	user.ID = result.InsertedID.(primitive.ObjectID)
	return err
}
