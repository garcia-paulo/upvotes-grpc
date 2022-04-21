package presentation

import (
	"net"

	"github.com/upvotes-grpc/garcia-paulo/infra/config"
	"github.com/upvotes-grpc/garcia-paulo/presentation/servers"
	"github.com/upvotes-grpc/garcia-paulo/proto/gen"
	"google.golang.org/grpc"
)

type Server struct {
	UserServer *servers.UserServer
	Config     *config.Config
}

func NewServer(userServer *servers.UserServer, config *config.Config) *Server {
	return &Server{
		UserServer: userServer,
		Config:     config,
	}
}

func (s *Server) RegisterServers(grpcServer *grpc.Server) {
	gen.RegisterUserServiceServer(grpcServer, s.UserServer)
}

func (s *Server) Run(grpcServer *grpc.Server) {

	s.RegisterServers(grpcServer)

	listener, err := net.Listen("tcp", s.Config.ServerPort)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
