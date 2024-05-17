package app

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func SetAuthInfoWithContext(ctx context.Context) (context.Context, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "get peer from context failed")
	}
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok || len(tlsAuth.State.PeerCertificates) == 0 || len(tlsAuth.State.PeerCertificates[0].DNSNames) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "get tls info from peer failed")
	}

	return SetAuthInfo(ctx, &AuthInfo{
		AppCode: tlsAuth.State.PeerCertificates[0].DNSNames[0],
	}), nil
}

type ModifiedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *ModifiedServerStream) Context() context.Context {
	return m.ctx
}

func UnaryAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, err = SetAuthInfoWithContext(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func StreamAuth() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := SetAuthInfoWithContext(ss.Context())
		if err != nil {
			return err
		}
		return handler(srv, &ModifiedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}
