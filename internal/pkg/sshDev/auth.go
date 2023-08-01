package sshDev

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TokenAuth(key string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("unauthorized")
		}

		token, ok := md["authorization"]
		if !ok || len(token) == 0 || token[0] != key {
			return nil, errors.New("unauthorized")
		}

		return handler(ctx, req)
	}
}
