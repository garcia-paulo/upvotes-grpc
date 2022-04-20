package servers

import (
	"context"

	"github.com/upvotes-grpc/garcia-paulo/proto/gen"
)

type UserServer struct {
	gen.UnimplementedUserServiceServer
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {}

func (s *UserServer) CreateUser(ctx context.Context, in *gen.UserRequest) (*gen.UserResponse, error) {
	return &gen.UserResponse{
		Id:       1,
		Username: in.Username,
	}, nil
}
