package presentation

import (
	"github.com/upvotes-grpc/garcia-paulo/presentation/servers"
	"github.com/upvotes-grpc/garcia-paulo/proto/gen"
	"google.golang.org/grpc"
)

func RegisterServers(grpcServer *grpc.Server) {
	gen.RegisterPostServiceServer(grpcServer, &servers.PostServer{})
	gen.RegisterUserServiceServer(grpcServer, &servers.UserServer{})
}
