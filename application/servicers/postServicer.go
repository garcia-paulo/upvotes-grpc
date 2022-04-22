package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostServicer struct {
	repository *repositories.PostRepository
}

func NewPostServicer(repository *repositories.PostRepository) *PostServicer {
	return &PostServicer{
		repository: repository,
	}
}

func (p *PostServicer) CreatePost(in *gen.PostRequest) (*gen.PostResponse, error) {
	post := models.NewPost(in)
	if err := post.Validate(); err != nil {
		return nil, err
	}

	err := p.repository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}

func (p *PostServicer) GetPosts() (*gen.ManyPostsResponse, error) {
	posts, err := p.repository.GetPosts()
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
	userId, err := primitive.ObjectIDFromHex(in.UserId)
	if err != nil {
		return nil, err
	}

	post, err := p.repository.GetPostById(postId)
	if err != nil {
		return nil, err
	}

	upvoteFound := false
	for _, upvote := range post.Upvotes {
		if upvote == userId {
			upvoteFound = true
			break
		}
	}

	if upvoteFound {
		err = p.repository.RemoveUpvote(post, userId)
	} else {
		err = p.repository.AddUpvote(post, userId)
	}
	if err != nil {
		return nil, err
	}

	return models.NewPostResponse(post), nil
}
