package main

import "google.golang.org/grpc"

func main() {
	server := InitializeServer()
	server.Run(grpc.NewServer())
}
