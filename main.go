package main

import (
	"github.com/garcia-paulo/upvotes-grpc/application/interceptors"
	"google.golang.org/grpc"
)

func main() {
	server := InitializeServer()
	server.Run(grpc.NewServer(grpc.UnaryInterceptor(interceptors.UnaryInterceptor)))
}
