// rpc server for applications
// mutual ssl verify is used

package app

import (
	"context"
	"errors"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/app/appProto"
)

type Server struct {
	appProto.UnimplementedAppServer
}

func (s *Server) GetTokenStatus(ctx context.Context, req *appProto.TokenRequest) (*appProto.TokenStatus, error) {
	err := redis.NewRecordedToken().NewStorePoint(req.Id).Get(ctx, nil)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &appProto.TokenStatus{
				Valid: false,
			}, nil
		}
	}
	return &appProto.TokenStatus{
		Valid: true,
	}, nil
}
