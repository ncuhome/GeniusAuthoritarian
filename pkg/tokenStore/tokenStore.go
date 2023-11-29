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

func (a TokenStore) CreateStorePoint(ctx context.Context, iat time.Time, valid time.Duration, claims interface{}) (uint64, error) {
	var value []byte
	var err error
	value, err = json.Marshal(&StorePointData{
		Iat:    iat.Unix(),
		Claims: claims,
	})
	if err != nil {
		return 0, err
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

func (a Point) parsePoint(iat time.Time, data []byte, claims interface{}) error {
	var parsedData StorePointData
	parsedData.Claims = claims
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		return err
	}
	if parsedData.Iat != iat.Unix() {
		return redis.Nil
	}
	return nil
}

func (a Point) GetAndDestroy(ctx context.Context, iat time.Time, claims interface{}) error {
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(iat, value, claims)
}

func (a Point) Get(ctx context.Context, iat time.Time, claims interface{}) error {
	value, err := a.s.client.Get(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(iat, value, claims)
}

func (a Point) Destroy(ctx context.Context) error {
	return a.s.client.Del(ctx, a.key).Err()
}
