package client

import "context"

type GrpcAuth struct {
	Token string
}

func (a *GrpcAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.Token}, nil
}

func (a *GrpcAuth) RequireTransportSecurity() bool {
	return true
}
