package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/application/token"
	"github.com/garcia-paulo/upvotes-grpc/server/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/repositories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServicer struct {
	repository *repositories.UserRepository
	tokenMaker *token.TokenMaker
}

func NewUserServicer(repository *repositories.UserRepository, tokenMaker *token.TokenMaker) *UserServicer {
	return &UserServicer{
		repository: repository,
		tokenMaker: tokenMaker,
	}
}

func (s *UserServicer) CreateUser(in *gen.UserRequest) (*gen.UserResponse, error) {
	user := models.NewUser(in)

	if err := user.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %s", err.Error())
	}

	hashedPassword, err := models.HashPassword(in.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password")
	}

	user.HashedPassword = hashedPassword

	err = s.repository.CreateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating user: %s", err.Error())
	}

	token, err := s.tokenMaker.CreateToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating token")
	}

	return models.NewUserResponse(user, token), nil
}

func (s *UserServicer) Login(in *gen.UserRequest) (*gen.UserResponse, error) {
	user, err := s.repository.FindUserByUsername(in.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user with username %s not found", in.Username)
	}

	if err := user.Authenticate(in.Password); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	token, err := s.tokenMaker.CreateToken(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating token")
	}

	return models.NewUserResponse(user, token), nil
}
