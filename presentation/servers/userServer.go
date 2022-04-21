package servers

import (
	"context"

	"github.com/upvotes-grpc/garcia-paulo/domain/models"
	"github.com/upvotes-grpc/garcia-paulo/infra/repositories"
	"github.com/upvotes-grpc/garcia-paulo/proto/gen"
)

type UserServer struct {
	gen.UnimplementedUserServiceServer
	repository *repositories.UserRepository
}

func NewUserServer(repository *repositories.UserRepository) *UserServer {
	return &UserServer{
		repository: repository,
	}
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {}

func (s *UserServer) CreateUser(ctx context.Context, in *gen.UserRequest) (*gen.UserResponse, error) {
	user := models.NewUser(in)

	if err := s.repository.CreateUser(user); err != nil {
		return nil, err
	}

	return models.NewUserResponse(user), nil
}
