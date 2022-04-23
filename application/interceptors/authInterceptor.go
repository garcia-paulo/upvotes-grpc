package interceptors

import (
	"context"
	"regexp"

	"github.com/garcia-paulo/upvotes-grpc/application/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	tokenMaker *token.TokenMaker
	protected  []string
}

func NewAuthInterceptor(tokenMaker *token.TokenMaker) *AuthInterceptor {
	return &AuthInterceptor{
		tokenMaker: tokenMaker,
		protected: []string{
			"/PostService/CreatePost",
			"/PostService/ToggleUpvote",
		},
	}
}

func (i *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newContext, err := i.Authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	if newContext != nil {
		return handler(newContext, req)
	}

	return handler(ctx, req)
}

func (i *AuthInterceptor) Authorize(ctx context.Context, method string) (context.Context, error) {
	protected := false
	for _, p := range i.protected {
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
