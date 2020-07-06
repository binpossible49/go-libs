package commonproto

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

type commonProtoServer struct {
}

func NewCommonProtoServer() CommonProtoAPIServer {
	return &commonProtoServer{}
}
func (c *commonProtoServer) Liveness(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &emptypb.Empty{}, nil

}
func (c *commonProtoServer) Readiness(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &emptypb.Empty{}, nil
}
