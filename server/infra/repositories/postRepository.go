package repositories

import (
	"context"

	"github.com/garcia-paulo/upvotes-grpc/server/domain/models"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostRepository struct {
	posts          *mongo.Collection
	userRepository *UserRepository
	ctx            context.Context
}

func NewPostRepository(database *database.Database, userRepository *UserRepository) *PostRepository {
	return &PostRepository{
		posts:          database.Database.Collection("posts"),
		userRepository: userRepository,
		ctx:            database.Ctx,
	}
}

func (p *PostRepository) CreatePost(post *models.Post) error {
	if _, err := p.userRepository.FindUserByUsername(post.Author); err != nil {
		return status.Errorf(codes.NotFound, "user with username %s not found", post.Author)
	}

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

func (p *PostRepository) RemoveUpvote(post *models.Post, username string) error {
	_, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$pull": bson.M{"upvotes": username}})
	if err != nil {
		return err
	}

	update, err := p.GetPostById(post.ID)
	post.Upvotes = update.Upvotes
	return err
}

func (p *PostRepository) AddUpvote(post *models.Post, username string) error {
	if _, err := p.userRepository.FindUserByUsername(username); err != nil {
		return status.Errorf(codes.NotFound, "user with username %s not found", username)
	}

	_, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$addToSet": bson.M{"upvotes": username}})
	post.Upvotes = append(post.Upvotes, username)
	return err
}
