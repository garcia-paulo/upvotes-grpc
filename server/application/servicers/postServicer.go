package servicers

import (
	"time"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/repositories"
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

func (p *PostServicer) GetPosts() (*gen.ManyPostsResponse, error) {
	posts, err := p.postRepository.GetPosts()
	if err != nil {
		return nil, err
	}

	return models.NewManyPostsResponse(posts), nil
}

func (p *PostServicer) CreatePost(in *gen.PostRequest, username string) (*gen.PostResponse, error) {
	post := models.NewPost(in, username)
	if err := post.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post: %s", err.Error())
	}

	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := p.postRepository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}

func (p *PostServicer) ToggleUpvote(in *gen.PostIdRequest, username string) (*gen.PostResponse, error) {
	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post id")
	}

	post, err := p.postRepository.GetPostById(postId)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return models.NewPostResponse(post), nil
}

func (p *PostServicer) DeletePost(in *gen.PostIdRequest, username string) (*gen.Message, error) {
	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post id")
	}

	err = p.postRepository.DeletePost(postId, username)
	if err != nil {
		return nil, err
	}

	return &gen.Message{
		Message: "post deleted",
	}, nil
}

func (p *PostServicer) UpdatePost(in *gen.PostUpdateRequest, username string) (*gen.PostResponse, error) {
	postId, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post id")
	}

	post, err := p.postRepository.GetPostById(postId)
	if err != nil {
		return nil, err
	}

	post.Title = in.Title
	post.Body = in.Body
	post.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := post.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post: %s", err.Error())
	}

	err = p.postRepository.UpdatePost(post, username)
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}
