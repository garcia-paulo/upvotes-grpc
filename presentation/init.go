package presentation

import (
	"net"

	"github.com/garcia-paulo/upvotes-grpc/application/interceptors"
	"github.com/garcia-paulo/upvotes-grpc/infra/config"
	"github.com/garcia-paulo/upvotes-grpc/presentation/servers"
	"github.com/garcia-paulo/upvotes-grpc/proto/gen"
	"google.golang.org/grpc"
)

type Server struct {
	UserServer      *servers.UserServer
	PostServer      *servers.PostServer
	AuthInterceptor *interceptors.AuthInterceptor
	config          *config.Config
}

func NewServer(userServer *servers.UserServer, config *config.Config, postServer *servers.PostServer, authInterceptor *interceptors.AuthInterceptor) *Server {
	return &Server{
		PostServer:      postServer,
		UserServer:      userServer,
		AuthInterceptor: authInterceptor,
		config:          config,
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
