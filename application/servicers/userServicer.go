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

	err := s.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return models.NewUserResponse(user), nil
}
