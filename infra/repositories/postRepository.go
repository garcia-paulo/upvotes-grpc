package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/infra/database"
	"go.mongodb.org/mongo-driver/bson"
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

func (p *PostRepository) GetPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	cursor, err := p.posts.Find(p.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(p.ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *PostRepository) GetPostById(id primitive.ObjectID) (*models.Post, error) {
	post := &models.Post{}
	err := p.posts.FindOne(p.ctx, bson.M{"_id": id}).Decode(post)
	return post, err
}

func (p *PostRepository) RemoveUpvote(post *models.Post, userId primitive.ObjectID) error {
	_, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$pull": bson.M{"upvotes": userId}})
	if err != nil {
		return err
	}

	update, err := p.GetPostById(post.ID)
	post.Upvotes = update.Upvotes
	return err
}

func (p *PostRepository) AddUpvote(post *models.Post, userId primitive.ObjectID) error {
	_, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$addToSet": bson.M{"upvotes": userId}})
	post.Upvotes = append(post.Upvotes, userId)
	return err
}
