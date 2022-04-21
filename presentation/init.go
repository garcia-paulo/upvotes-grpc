package presentation

import (
	"net"

	"github.com/garcia-paulo/upvotes-grpc/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/presentation/servers"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"google.golang.org/grpc"
)

type Server struct {
	userServer *servers.UserServer
	postServer *servers.PostServer
	config     *config.Config
}

func NewServer(userServer *servers.UserServer, config *config.Config, postServer *servers.PostServer) *Server {
	return &Server{
		postServer: postServer,
		userServer: userServer,
		config:     config,
	}
}

func (s *Server) RegisterServers(grpcServer *grpc.Server) {
	gen.RegisterUserServiceServer(grpcServer, s.userServer)
	gen.RegisterPostServiceServer(grpcServer, s.postServer)
}

func (s *Server) Run(grpcServer *grpc.Server) {

	s.RegisterServers(grpcServer)

	listener, err := net.Listen("tcp", s.config.ServerPort)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
