package servicers

import (
	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/repositories"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
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
