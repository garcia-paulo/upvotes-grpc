package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
)

type UserServicer struct {
	repository *repositories.UserRepository
}

func NewUserServicer(repository *repositories.UserRepository) *UserServicer {
	return &UserServicer{
		repository: repository,
	}
}

func (s *UserServicer) CreateUser(in *gen.UserRequest) (*gen.UserResponse, error) {
	user := models.NewUser(in)

	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := models.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = hashedPassword

	err = s.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return models.NewUserResponse(user), nil
}

func (s *UserServicer) Login(in *gen.UserRequest) (*gen.UserResponse, error) {
	user, err := s.repository.FindUserByUsername(in.Username)
	if err != nil {
		return nil, err
	}

	if err := user.Authenticate(in.Password); err != nil {
		return nil, err
	}

	return models.NewUserResponse(user), nil
}
