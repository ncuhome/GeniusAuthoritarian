package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type UserOperationIDChannel struct {
	key string
}

type UserOperationIDInfo struct {
	UID         uint   `json:"uid"`
	OperationID uint64 `json:"operationID"`
}

func NewUserOperationIDChannel() UserOperationIDChannel {
	return UserOperationIDChannel{
		key: keyUserJwt.String() + "sub",
	}
}

func (channel UserOperationIDChannel) Publish(ctx context.Context, info UserOperationIDInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return Client.Publish(ctx, channel.key, data).Err()
}

func (channel UserOperationIDChannel) Subscribe(ctx context.Context) *redis.PubSub {
	return Client.Subscribe(ctx, channel.key)
}

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
		uid:     uid,
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
	uid     uint
}

func (u UserJwtOperator) Create(ctx context.Context) error {
	err := NewUserOperationIDChannel().Publish(ctx, UserOperationIDInfo{
		UID:         u.uid,
		OperationID: 1,
	})
	if err != nil {
		return err
	}
	return Client.HSet(ctx, u.key, u.userKey, 1).Err()
}

func (u UserJwtOperator) Del(ctx context.Context) error {
	err := NewUserOperationIDChannel().Publish(ctx, UserOperationIDInfo{
		UID:         u.uid,
		OperationID: 0,
	})
	if err != nil {
		return err
	}
	return Client.HDel(ctx, u.key, u.userKey).Err()
}

func (u UserJwtOperator) Exist(ctx context.Context) (bool, error) {
	return Client.HExists(ctx, u.key, u.userKey).Result()
}

func (u UserJwtOperator) GetOperateID(ctx context.Context) (uint64, error) {
	return Client.HGet(ctx, u.key, u.userKey).Uint64()
}

func (u UserJwtOperator) ChangeOperateID(ctx context.Context) (uint64, error) {
	newOperationID, err := Client.HIncrBy(ctx, u.key, u.userKey, 1).Uint64()
	if err != nil {
		return 0, err
	}
	err = NewUserOperationIDChannel().Publish(ctx, UserOperationIDInfo{
		UID:         u.uid,
		OperationID: newOperationID,
	})
	if err != nil {
		return 0, err
	}
	return newOperationID, err
}

func (u UserJwtOperator) CheckOperateID(ctx context.Context, oid uint64) (bool, error) {
	id, err := u.GetOperateID(ctx)
	return id == oid, err
}
