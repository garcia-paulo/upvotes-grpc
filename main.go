package main

import "google.golang.org/grpc"

func main() {
	grpcServer := grpc.NewServer()
	servers := InitializeServer()
	servers.RegisterServers(grpcServer)
	servers.Run(grpcServer)
}
