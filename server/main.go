package main

import (
	"google.golang.org/grpc"
)

func main() {
	server := InitializeServer()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.UnaryInterceptor.Unary))
	server.RegisterServers(grpcServer)
	server.Run(grpcServer)
}
