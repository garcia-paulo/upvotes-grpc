package servers

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/application/servicers"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
)

type PostServer struct {
	gen.UnimplementedPostServiceServer
	servicer *servicers.PostServicer
}

func NewPostServer(servicer *servicers.PostServicer) *PostServer {
	return &PostServer{
		servicer: servicer,
	}
}

func (s *PostServer) mustEmbedUnimplementedPostServiceServer() {}

func (s *PostServer) CreatePost(ctx context.Context, in *gen.PostRequest) (*gen.PostResponse, error) {
	response, err := s.servicer.CreatePost(in)
	if err != nil {
		return nil, err
	}

	return response, nil
}
