package middlewares

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

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
			log.Tracef("RPC %s [%s] %s - success", info.FullMethod, remoteAddr, subTimeStr)
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
