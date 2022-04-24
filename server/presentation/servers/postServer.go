package servers

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/application/servicers"
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

func (s *PostServer) GetPosts(ctx context.Context, in *gen.Empty) (*gen.ManyPostsResponse, error) {
	response, err := s.servicer.GetPosts()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PostServer) CreatePost(ctx context.Context, in *gen.PostRequest) (*gen.PostResponse, error) {
	response, err := s.servicer.CreatePost(in, ctx.Value("username").(string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PostServer) ToggleUpvote(ctx context.Context, in *gen.PostIdRequest) (*gen.PostResponse, error) {
	response, err := s.servicer.ToggleUpvote(in, ctx.Value("username").(string))
	if err != nil {
		return nil, err
	}

	return response, nil
}
