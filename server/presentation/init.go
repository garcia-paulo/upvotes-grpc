package presentation

import (
	"net"

	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"github.com/garcia-paulo/upvotes-grpc/server/application/interceptors"
	"github.com/garcia-paulo/upvotes-grpc/server/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/server/presentation/servers"
	"google.golang.org/grpc"
)

type Server struct {
	UserServer       *servers.UserServer
	PostServer       *servers.PostServer
	UnaryInterceptor *interceptors.UnaryInterceptor
	config           *config.Config
}

func NewServer(userServer *servers.UserServer, config *config.Config, postServer *servers.PostServer, unaryInterceptor *interceptors.UnaryInterceptor) *Server {
	return &Server{
		PostServer:       postServer,
		UserServer:       userServer,
		UnaryInterceptor: unaryInterceptor,
		config:           config,
	}
}

func (s *Server) RegisterServers(grpcServer *grpc.Server) {
	gen.RegisterUserServiceServer(grpcServer, s.UserServer)
	gen.RegisterPostServiceServer(grpcServer, s.PostServer)
}

func (s *Server) Run(grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", s.config.ServerPort)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
