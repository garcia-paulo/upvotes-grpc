package models

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/validator.v2"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Username       string             `bson:"name" validate:"regexp=^[a-zA-Z0-9.-_]+$,min=4,max=12"`
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

func (u *User) Validate() error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}
