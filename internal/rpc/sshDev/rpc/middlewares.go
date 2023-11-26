package rpc

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

func verifyToken(ctx context.Context, key string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "token not found")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 || token[0] != key {
		return status.Error(codes.Unauthenticated, "insufficient permissions")
	}

	return nil
}

func UnaryTokenAuth(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = verifyToken(ctx, key); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func StreamTokenAuth(key string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := verifyToken(ss.Context(), key); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		remote, _ := peer.FromContext(ctx)
		remoteAddr := remote.Addr.String()

		startAt := time.Now()

		resp, err = handler(ctx, req)

		var subTimeStr string
		subTime := time.Now().Sub(startAt)
		if subTime > time.Millisecond {
			milliSeconds := subTime.Milliseconds()
			subTimeStr = fmt.Sprint(milliSeconds%1000) + "ms"
			milliSeconds = milliSeconds / 1000
			if milliSeconds > 0 {
				subTimeStr = fmt.Sprint(milliSeconds) + "s" + subTimeStr
			}
		} else {
			subTimeStr = fmt.Sprint(subTime.Microseconds(), "µs")
		}

		if err == nil {
			log.Infof("RPC %s [%s] %s - success", info.FullMethod, remoteAddr, subTimeStr)
		} else {
			log.Warnf("RPC %s [%s] %s - with error: %s", info.FullMethod, remoteAddr, subTimeStr, err)
		}

		return resp, err
	}
}

func StreamLogger() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		remote, _ := peer.FromContext(ss.Context())
		remoteAddr := remote.Addr.String()
		log.Infof("RPC STREAM %s [%s] - start", info.FullMethod, remoteAddr)

		err := handler(srv, ss)
		if err != nil {
			log.Warnf("RPC STREAM %s [%s] - end with error: %s", info.FullMethod, remoteAddr, err)
		} else {
			log.Infof("RPC STREAM %s [%s] - complete", info.FullMethod, remoteAddr)
		}
		return err
	}
}
