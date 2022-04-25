package models

import (
	"time"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/validator.v2"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title" validate:"min=6,max=18"`
	Body      string             `bson:"body" validate:"max=280,nonzero"`
	Upvotes   []string           `bson:"upvotes,omitempty"`
	Author    string             `bson:"author" validate:"nonzero"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt"`
}

func NewPost(request *gen.PostRequest, username string) *Post {
	return &Post{
		ID:        primitive.NewObjectID(),
		Title:     request.Title,
		Body:      request.Body,
		Author:    username,
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
}

func NewPostResponse(post *Post) *gen.PostResponse {
	response := &gen.PostResponse{
		Id:        post.ID.Hex(),
		Title:     post.Title,
		Body:      post.Body,
		Author:    post.Author,
		CreatedAt: post.UpdatedAt.Time().Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Time().Format(time.RFC3339),
	}

	response.Upvotes = append(response.Upvotes, post.Upvotes...)
	return response
}

func NewManyPostsResponse(posts []*Post) *gen.ManyPostsResponse {
	response := &gen.ManyPostsResponse{}
	for _, post := range posts {
		response.Posts = append(response.Posts, NewPostResponse(post))
	}
	return response
}

func (p *Post) Validate() error {
	if err := validator.Validate(p); err != nil {
		return err
	}
	return nil
}
