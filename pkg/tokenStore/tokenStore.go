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

type Node struct {
	// this is a unique id for each node process, it
	// will be reallocated every day.
	// ID must be smaller than 100.
	ID uint64

	client          *redis.Client
	keyNodeIDPrefix string

	// use for refresh fields.
	Lock       *sync.Mutex
	IDTimeMark uint64
	TokenID    *atomic.Uint64
}

func (node *Node) keyNodeID(timeMark uint64) string {
	return fmt.Sprintf("%s-%d", node.keyNodeIDPrefix, timeMark)
}

func (node *Node) currentTimeMark() uint64 {
	return uint64(time.Now().YearDay())
}

func (node *Node) GenID(ctx context.Context) (uint64, error) {
	currentTimeMark := node.currentTimeMark()
	if node.IDTimeMark != currentTimeMark {
		node.Lock.Lock()
		if node.IDTimeMark == currentTimeMark {
			node.Lock.Unlock()
			return node.GenID(ctx)
		}
		defer node.Lock.Unlock()
		newNodeID, err := node.client.Incr(ctx, node.keyNodeID(currentTimeMark)).Uint64()
		if err != nil {
			return 0, err
		}
		node.ID = newNodeID % 100
		node.TokenID.Store(0)
		node.IDTimeMark = currentTimeMark
	}
	tokenID := node.TokenID.Add(1)
	tokenID = (tokenID << 5) + (node.ID << 3) + node.IDTimeMark
	return tokenID, nil
}

func NewTokenStoreFactory[C any](Client *redis.Client, keyPrefix string) func() TokenStore[C] {
	node := Node{
		client:          Client,
		keyNodeIDPrefix: keyPrefix + "id",
		Lock:            &sync.Mutex{},
		TokenID:         &atomic.Uint64{},
	}
	return func() TokenStore[C] {
		return TokenStore[C]{
			client:    Client,
			keyPrefix: keyPrefix,
			node:      &node,
		}
	}
}

type TokenStore[C any] struct {
	client *redis.Client
	// token 有效校验的 key 前缀
	keyPrefix string

	// node should be set with static variable
	node *Node
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
