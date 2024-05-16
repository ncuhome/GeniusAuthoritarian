package tokenStore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"sync/atomic"
	"time"
)

func NewTokenStoreFactory[C any](keyPrefix string, getClient func() *redis.Client) func() TokenStore[C] {
	node := &Node{
		keyNodeIDPrefix: keyPrefix + "id",
		Lock:            &sync.Mutex{},
		TokenID:         &atomic.Uint64{},
	}
	return func() TokenStore[C] {
		client := getClient()
		return TokenStore[C]{
			client:    getClient(),
			keyPrefix: keyPrefix,
			node:      node.WithClient(client),
		}
	}
}

type TokenStore[C any] struct {
	client *redis.Client
	// token 有效校验的 key 前缀
	keyPrefix string

	// node should be set with static variable
	node NodeWithClient
}

func (store TokenStore[C]) genKey(id uint64) string {
	return store.keyPrefix + fmt.Sprint(id)
}

func (store TokenStore[C]) CreateStorePointWithID(ctx context.Context, id uint64, valid time.Duration, claims *C) error {
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
	return store.client.Set(ctx, store.genKey(id), value, valid).Err()
}

func (store TokenStore[C]) CreateStorePoint(ctx context.Context, valid time.Duration, claims *C) (uint64, error) {
	id, err := store.node.GenID(ctx)
	if err != nil {
		return 0, err
	}
	return id, store.CreateStorePointWithID(ctx, id, valid, claims)
}

func (store TokenStore[C]) MPointGet(ctx context.Context, ids ...uint64) ([]interface{}, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = store.genKey(id)
	}
	return store.client.MGet(ctx, keys...).Result()
}

func (store TokenStore[C]) NewStorePoint(id uint64) Point[C] {
	return Point[C]{
		s:   store,
		key: store.genKey(id),
	}
}

type Point[C any] struct {
	s   TokenStore[C]
	key string
}

func (point Point[C]) parsePoint(data []byte, claims *C) error {
	if claims != nil {
		return json.Unmarshal(data, claims)
	}
	return nil
}

func (point Point[C]) GetAndDestroy(ctx context.Context, claims *C) error {
	value, err := point.s.client.GetDel(ctx, point.key).Bytes()
	if err != nil {
		return err
	}

	return point.parsePoint(value, claims)
}

func (point Point[C]) Get(ctx context.Context, claims *C) error {
	value, err := point.s.client.Get(ctx, point.key).Bytes()
	if err != nil {
		return err
	}

	return point.parsePoint(value, claims)
}

func (point Point[C]) Destroy(ctx context.Context) error {
	return point.s.client.Del(ctx, point.key).Err()
}
