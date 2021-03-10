package middleware

import (
	"context"

	"github.com/xulichen/halfway/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

// Validate return a server interceptor validate incoming request per RPC call.
func Validate() grpc.UnaryServerInterceptor {
	validate := utils.NewValidator()
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = validate.Validate(req); err != nil {
			err = status.Error(codes.InvalidArgument, err.Error())
			return
		}
		resp, err = handler(ctx, req)
		return
	}
}
