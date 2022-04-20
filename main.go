package main

import (
	"net"

	"github.com/upvotes-grpc/garcia-paulo/presentation"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	presentation.RegisterServers(grpcServer)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
