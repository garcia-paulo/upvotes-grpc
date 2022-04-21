package models

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Username       string             `bson:"name"`
	HashedPassword string             `bson:"hashed_password"`
}

func NewUser(request *gen.UserRequest) *User {
	return &User{
		ID:             primitive.NewObjectID(),
		Username:       request.Username,
		HashedPassword: request.Password,
	}
}

func NewUserResponse(user *User) *gen.UserResponse {
	return &gen.UserResponse{
		Id:       user.ID.Hex(),
		Username: user.Username,
	}
}
