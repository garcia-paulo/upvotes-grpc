package main

import (
	"google.golang.org/grpc"
)

func main() {
	server := InitializeServer()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.AuthInterceptor.UnaryInterceptor))
	server.RegisterServers(grpcServer)
	server.Run(grpcServer)
}
