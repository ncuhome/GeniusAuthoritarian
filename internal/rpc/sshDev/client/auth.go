package client

import "context"

type GrpcAuth struct {
	Token string
}

func (a *GrpcAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.Token}, nil
}

func (a *GrpcAuth) RequireTransportSecurity() bool {
	return true
}
