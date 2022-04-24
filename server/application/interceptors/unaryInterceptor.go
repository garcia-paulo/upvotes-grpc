package interceptors

import (
	"context"
	"regexp"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/application/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UnaryInterceptor struct {
	tokenMaker    *token.TokenMaker
	authProtected []string
}

func NewUnaryInterceptor(tokenMaker *token.TokenMaker) *UnaryInterceptor {
	return &UnaryInterceptor{
		tokenMaker: tokenMaker,
		authProtected: []string{
			"/PostService/CreatePost",
			"/PostService/ToggleUpvote",
		},
	}
}

func (i *UnaryInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/UserService/CreateUser" {
		err := i.invalidUserInterceptor(req.(*gen.UserRequest))
		if err != nil {
			return nil, err
		}
	}

	newContext, err := i.authInterceptor(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	if newContext != nil {
		return handler(newContext, req)
	}

	return handler(ctx, req)
}

func (i *UnaryInterceptor) invalidUserInterceptor(req *gen.UserRequest) error {
	invalidUsernames := []string{
		"test.user",
		"test.user2",
	}

	for _, u := range invalidUsernames {
		if u == req.Username {
			return status.Errorf(codes.Unauthenticated, "invalid username")
		}
	}

	return nil
}

func (i *UnaryInterceptor) authInterceptor(ctx context.Context, method string) (context.Context, error) {
	protected := false
	for _, p := range i.authProtected {
		r, _ := regexp.MatchString(p, method)
		if r {
			protected = true
			break
		}
	}

	if !protected {
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	token := values[0]
	payload, err := i.tokenMaker.VerifyToken(token)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	newContext := context.WithValue(ctx, "username", payload.Username)

	return newContext, nil
}
