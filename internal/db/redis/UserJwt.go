package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// NewUserJwt user jwt Operate ID hash 表
func NewUserJwt() UserJwt {
	return UserJwt{
		key: keyUserJwt.String(),
	}
}

type UserJwt struct {
	key string
}

func (u UserJwt) Empty(ctx context.Context) (bool, error) {
	length, err := Client.HLen(ctx, u.key).Result()
	if errors.Is(err, redis.Nil) {
		return true, nil
	}
	return length == 0, err
}

func (u UserJwt) NewOperator(uid uint) UserJwtOperator {
	return UserJwtOperator{
		key:     u.key,
		userKey: fmt.Sprint(uid),
	}
}

type UserJwtOperator struct {
	key     string
	userKey string
}

func (u UserJwtOperator) Create(ctx context.Context) error {
	return Client.HSet(ctx, u.key, u.userKey, 0).Err()
}

func (u UserJwtOperator) Del(ctx context.Context) error {
	return Client.HDel(ctx, u.key, u.userKey).Err()
}

func (u UserJwtOperator) Exist(ctx context.Context) (bool, error) {
	return Client.HExists(ctx, u.key, u.userKey).Result()
}

func (u UserJwtOperator) GetOperateID(ctx context.Context) (uint64, error) {
	return Client.HGet(ctx, u.key, u.userKey).Uint64()
}

func (u UserJwtOperator) ChangeOperateID(ctx context.Context) (uint64, error) {
	return Client.HIncrBy(ctx, u.key, u.userKey, 1).Uint64()
}

func (u UserJwtOperator) CheckOperateID(ctx context.Context, oid uint64) (bool, error) {
	id, err := u.GetOperateID(ctx)
	return id == oid, err
}
