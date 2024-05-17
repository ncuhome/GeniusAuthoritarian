package middlewares

import (
	"context"
	"github.com/Mmx233/tool"
	"google.golang.org/grpc"
)

func UnaryRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer tool.Recover()
		return handler(ctx, req)
	}
}

func StreamRecovery() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer tool.Recover()
		return handler(srv, ss)
	}
}
