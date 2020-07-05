package commonproto

import (
	context "context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc "google.golang.org/grpc"
)

func RegisterCommonProto(server *grpc.Server) {
	commonAPIServer := NewCommonProtoServer()
	RegisterCommonProtoAPIServer(server, commonAPIServer)
}

func GRPCGatewayHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return RegisterCommonProtoAPIHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
