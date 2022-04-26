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
	if err != nil {
		return status.Errorf(codes.Internal, "error creating post: %s", err.Error())
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (p *PostRepository) GetPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	cursor, err := p.posts.Find(p.ctx, bson.D{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting posts: %s", err.Error())
	}

	for cursor.Next(p.ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, status.Errorf(codes.Internal, "error decoding post: %s", err.Error())
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *PostRepository) GetPostById(id primitive.ObjectID) (*models.Post, error) {
	post := &models.Post{}
	err := p.posts.FindOne(p.ctx, bson.M{"_id": id}).Decode(post)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}

	return post, nil
}

func (p *PostRepository) RemoveUpvote(post *models.Post, username string) error {
	_, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$pull": bson.M{"upvotes": username}})
	if err != nil {
		return status.Errorf(codes.Internal, "error removing upvote: %s", err.Error())
	}

	update, err := p.GetPostById(post.ID)
	if err != nil {
		return status.Errorf(codes.Internal, "error when retrieving post: %s", err.Error())
	}
	post.Upvotes = update.Upvotes
	return nil
}

func (p *PostRepository) AddUpvote(post *models.Post, username string) error {
	if _, err := p.userRepository.FindUserByUsername(username); err != nil {
		return status.Errorf(codes.NotFound, "user with username %s not found", username)
	}

	if _, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$addToSet": bson.M{"upvotes": username}}); err != nil {
		return status.Errorf(codes.Internal, "error adding upvote: %s", err.Error())
	}
	post.Upvotes = append(post.Upvotes, username)
	return nil
}

func (p *PostRepository) DeletePost(id primitive.ObjectID, username string) error {
	post, err := p.GetPostById(id)
	if err != nil {
		return status.Errorf(codes.NotFound, "post not found")
	}

	if post.Author != username {
		return status.Errorf(codes.PermissionDenied, "user %s is not the author of the post", username)
	}

	_, err = p.posts.DeleteOne(p.ctx, bson.M{"_id": id})
	if err != nil {
		return status.Errorf(codes.NotFound, "post not found")
	}

	return nil
}

func (p *PostRepository) UpdatePost(post *models.Post, username string) error {
	if post.Author != username {
		return status.Errorf(codes.PermissionDenied, "user %s is not the author of the post", username)
	}

	if _, err := p.posts.UpdateOne(p.ctx, bson.M{"_id": post.ID}, bson.M{"$set": post}); err != nil {
		return status.Errorf(codes.Internal, "error updating post: %s", err.Error())
	}

	return nil
}
