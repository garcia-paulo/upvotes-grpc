package models

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/validator.v2"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Username       string             `bson:"username" validate:"regexp=^[a-zA-Z0-9.-_]+$,min=4,max=12"`
	HashedPassword string             `bson:"hashed_password" validate:"regexp=^[a-zA-Z0-9.-_@#]+$,min=6,max=18"`
}

func NewUser(request *gen.UserRequest) *User {
	return &User{
		ID:             primitive.NewObjectID(),
		Username:       request.Username,
		HashedPassword: request.Password,
	}
}

func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid password")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", status.Errorf(codes.Internal, "error hashing password")
	}
	return string(hashedPassword), nil
}

func NewUserResponse(user *User, token string) *gen.UserResponse {
	return &gen.UserResponse{
		Username: user.Username,
		Token:    token,
	}
}

func (u *User) Validate() error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}
