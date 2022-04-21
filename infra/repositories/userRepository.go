package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (u *UserRepository) CreateUser(user *models.User) error {
	result, err := u.users.InsertOne(u.ctx, user)
	user.ID = result.InsertedID.(primitive.ObjectID)
	return err
}
