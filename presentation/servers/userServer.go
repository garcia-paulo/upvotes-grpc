package servers

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/application/servicers"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
)

type UserServer struct {
	gen.UnimplementedUserServiceServer
	servicer *servicers.UserServicer
}

func NewUserServer(servicer *servicers.UserServicer) *UserServer {
	s := &UserServer{
		servicer: servicer,
	}
	s.mustEmbedUnimplementedUserServiceServer()
	return s
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {}

func (s *UserServer) CreateUser(ctx context.Context, in *gen.UserRequest) (*gen.UserResponse, error) {
	response, err := s.servicer.CreateUser(in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserServer) Login(ctx context.Context, in *gen.UserRequest) (*gen.UserResponse, error) {
	response, err := s.servicer.Login(in)
	if err != nil {
		return nil, err
	}

	return response, nil
}
