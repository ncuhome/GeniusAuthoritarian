package jwtVerify

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
)

func CheckUserClaims(ctx context.Context, claims jwtClaims.ClaimsUser) (bool, error) {
	return redis.NewUserJwt().NewOperator(claims.GetUID()).CheckOperateID(ctx, claims.GetUserOperateID())
}
