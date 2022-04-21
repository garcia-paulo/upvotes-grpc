package models

import (
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title"`
	Body  string             `bson:"body"`
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
