package models

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/validator.v2"
)

type Post struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title" validate:"min=6,max=18"`
	Body  string             `bson:"body" validate:"max=280,nonzero"`
}

func NewPost(request *gen.PostRequest) *Post {
	return &Post{
		ID:    primitive.NewObjectID(),
		Title: request.Title,
		Body:  request.Body,
	}
}

func NewPostResponse(post *Post) *gen.PostResponse {
	return &gen.PostResponse{
		Id:    post.ID.Hex(),
		Title: post.Title,
		Body:  post.Body,
	}
}

func (p *Post) Validate() error {
	if err := validator.Validate(p); err != nil {
		return err
	}
	return nil
}
