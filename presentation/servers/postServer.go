package servers

import (
	"context"

	"github.com/upvotes-grpc/garcia-paulo/proto/gen"
)

type PostServer struct {
	gen.UnimplementedPostServiceServer
}

func (s *PostServer) mustEmbedUnimplementedPostServiceServer() {}

func (s *PostServer) CreatePost(ctx context.Context, in *gen.PostRequest) (*gen.PostResponse, error) {
	return &gen.PostResponse{
		Title: "Hello World",
		Text:  "Hello World",
	}, nil
}
