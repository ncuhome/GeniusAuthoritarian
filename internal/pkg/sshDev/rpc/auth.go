package rpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func verifyToken(ctx context.Context, key string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("unauthorized")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 || token[0] != key {
		return errors.New("unauthorized")
	}

	return nil
}

func TokenAuthUnary(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = verifyToken(ctx, key); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func TokenAuthStream(key string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := verifyToken(ss.Context(), key); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}
