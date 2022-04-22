package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostServicer struct {
	postRepository *repositories.PostRepository
	userRepository *repositories.UserRepository
}

func NewPostServicer(postRepository *repositories.PostRepository, userRepository *repositories.UserRepository) *PostServicer {
	return &PostServicer{
		postRepository: postRepository,
		userRepository: userRepository,
	}
}

func (p *PostServicer) CreatePost(in *gen.PostRequest) (*gen.PostResponse, error) {
	post := models.NewPost(in)
	if err := post.Validate(); err != nil {
		return nil, err
	}

	err := p.postRepository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}

func (p *PostServicer) GetPosts() (*gen.ManyPostsResponse, error) {
	posts, err := p.postRepository.GetPosts()
	if err != nil {
		return nil, err
	}

	return models.NewManyPostsResponse(posts), nil
}

func (p *PostServicer) ToggleUpvote(in *gen.PostUserIdRequest) (*gen.PostResponse, error) {
	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, err
	}

	post, err := p.postRepository.GetPostById(postId)
	if err != nil {
		return nil, err
	}

	if _, err := p.userRepository.FindUserByUsername(in.Username); err != nil {
		return nil, err
	}

	upvoteFound := false
	for _, upvote := range post.Upvotes {
		if upvote == in.Username {
			upvoteFound = true
			break
		}
	}

	if upvoteFound {
		err = p.postRepository.RemoveUpvote(post, in.Username)
	} else {
		err = p.postRepository.AddUpvote(post, in.Username)
	}
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}
