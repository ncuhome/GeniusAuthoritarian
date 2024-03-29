package tokenStore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewTokenStore[C any](Client *redis.Client, keyPrefix string) TokenStore[C] {
	return TokenStore[C]{
		client:    Client,
		keyPrefix: keyPrefix,
		keyID:     keyPrefix + "id",
	}
}

type TokenStore[C any] struct {
	client *redis.Client
	// token 有效校验的 key 前缀
	keyPrefix string
	// redis ID 字段 key，用于给 token 分配不一样的 ID
	keyID string
}

func (a TokenStore[C]) genKey(id uint64) string {
	return a.keyPrefix + fmt.Sprint(id)
}

func (a TokenStore[C]) CreateStorePointWithID(ctx context.Context, id uint64, valid time.Duration, claims *C) error {
	var value []byte
	var err error
	if claims != nil {
		value, err = json.Marshal(claims)
		if err != nil {
			return err
		}
	} else {
		value = []byte{'1'}
	}
	return a.client.Set(ctx, a.genKey(id), value, valid).Err()
}

func (a TokenStore[C]) CreateStorePoint(ctx context.Context, valid time.Duration, claims *C) (uint64, error) {
	id, err := a.client.Incr(ctx, a.keyID).Uint64()
	if err != nil {
		return 0, err
	}
	return id, a.CreateStorePointWithID(ctx, id, valid, claims)
}

func (a TokenStore[C]) NewStorePoint(id uint64) Point[C] {
	return Point[C]{
		s:   a,
		key: a.genKey(id),
	}
}

type Point[C any] struct {
	s   TokenStore[C]
	key string
}

func (a Point[C]) parsePoint(data []byte, claims *C) error {
	if claims != nil {
		return json.Unmarshal(data, claims)
	}
	return nil
}

func (a Point[C]) GetAndDestroy(ctx context.Context, claims *C) error {
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(value, claims)
}

func (a Point[C]) Get(ctx context.Context, claims *C) error {
	value, err := a.s.client.Get(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(value, claims)
}

func (a Point[C]) Destroy(ctx context.Context) error {
	return a.s.client.Del(ctx, a.key).Err()
}
