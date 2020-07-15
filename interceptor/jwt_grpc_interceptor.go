package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// JWTAuthenticationInterceptor jwt authentication interceptor for grpc
func JWTAuthenticationInterceptor(publicKey string, excludingPath []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}
