package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
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

func (u UserJwt) GetOperationTable(ctx context.Context) (map[uint]uint64, error) {
	values, err := Client.HGetAll(ctx, u.key).Result()
	if err != nil {
		return nil, err
	}
	var result = make(map[uint]uint64, len(values))
	for uidStr, operationIDStr := range values {
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			return nil, err
		}
		operationID, err := strconv.ParseUint(operationIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		result[uint(uid)] = operationID
	}
	return result, nil
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
