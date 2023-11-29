package tokenStore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewTokenStore(Client *redis.Client, keyPrefix string) TokenStore {
	return TokenStore{
		client:    Client,
		keyPrefix: keyPrefix,
		keyID:     keyPrefix + "id",
	}
}

type TokenStore struct {
	client *redis.Client
	// token 有效校验的 key 前缀
	keyPrefix string
	// redis ID 字段 key，用于给 token 分配不一样的 ID
	keyID string
}

func (a TokenStore) genKey(id uint64) string {
	return a.keyPrefix + fmt.Sprint(id)
}

func (a TokenStore) CreateStorePoint(ctx context.Context, valid time.Duration, claims interface{}) (uint64, error) {
	var value []byte
	var err error
	if claims != nil {
		value, err = json.Marshal(claims)
		if err != nil {
			return 0, err
		}
	} else {
		value = []byte{'1'}
	}

	id, err := a.client.Incr(ctx, a.keyID).Uint64()
	if err != nil {
		return 0, err
	}

	return id, a.client.Set(ctx, a.genKey(id), value, valid).Err()
}

func (a TokenStore) NewStorePoint(id uint64) Point {
	return Point{
		s:   a,
		key: a.genKey(id),
	}
}

type Point struct {
	s   TokenStore
	key string
}

func (a Point) parsePoint(data []byte, claims interface{}) error {
	if claims != nil {
		return json.Unmarshal(data, claims)
	}
	return nil
}

func (a Point) GetAndDestroy(ctx context.Context, claims interface{}) error {
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(value, claims)
}

func (a Point) Get(ctx context.Context, claims interface{}) error {
	value, err := a.s.client.Get(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(value, claims)
}

func (a Point) Destroy(ctx context.Context) error {
	return a.s.client.Del(ctx, a.key).Err()
}
