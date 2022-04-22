package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("unaryInterceptor")
	return handler(ctx, req)
}
