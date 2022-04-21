package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	posts *mongo.Collection
	ctx   context.Context
}

func NewPostRepository(database *database.Database) *PostRepository {
	return &PostRepository{
		posts: database.Database.Collection("posts"),
		ctx:   database.Ctx,
	}
}

func (p *PostRepository) CreatePost(post *models.Post) error {
	result, err := p.posts.InsertOne(p.ctx, post)
	post.ID = result.InsertedID.(primitive.ObjectID)
	return err
}
