package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServicer struct {
	postRepository *repositories.PostRepository
}

func NewPostServicer(postRepository *repositories.PostRepository) *PostServicer {
	return &PostServicer{
		postRepository: postRepository,
	}
}

func (p *PostServicer) CreatePost(in *gen.PostRequest, username string) (*gen.PostResponse, error) {
	post := models.NewPost(in, username)
	if err := post.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post: %s", err.Error())
	}

	err := p.postRepository.CreatePost(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating post: %s", err.Error())
	}

	return models.NewPostResponse(post), nil
}

func (p *PostServicer) GetPosts() (*gen.ManyPostsResponse, error) {
	posts, err := p.postRepository.GetPosts()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting posts")
	}

	return models.NewManyPostsResponse(posts), nil
}

func (p *PostServicer) ToggleUpvote(in *gen.PostIdRequest, username string) (*gen.PostResponse, error) {
	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post id")
	}

	post, err := p.postRepository.GetPostById(postId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}

	upvoteFound := false
	for _, upvote := range post.Upvotes {
		if upvote == username {
			upvoteFound = true
			break
		}
	}

	if upvoteFound {
		err = p.postRepository.RemoveUpvote(post, username)
	} else {
		err = p.postRepository.AddUpvote(post, username)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error toggling upvote: %s", err.Error())
	}

	return models.NewPostResponse(post), nil
}
