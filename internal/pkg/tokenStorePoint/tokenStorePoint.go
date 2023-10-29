package tokenStorePoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
	"time"
)

func NewTokenStore(Client *redis.Client, idAtom *atomic.Uint64, keyPrefix string, iat time.Time) TokenStore {
	return TokenStore{
		client:    Client,
		keyPrefix: keyPrefix,
		iat:       iat.Unix(),
		idAtom:    idAtom,
	}
}

type TokenStore struct {
	client    *redis.Client
	keyPrefix string
	iat       int64
	idAtom    *atomic.Uint64
}

func (a TokenStore) genKey(id uint64) string {
	return a.keyPrefix + fmt.Sprint(id)
}

func (a TokenStore) CreateStorePoint(ctx context.Context, valid time.Duration, claims interface{}) (uint64, error) {
	value, err := json.Marshal(&StorePointData{
		Iat:    a.iat,
		Claims: claims,
	})
	if err != nil {
		return 0, err
	}

	id := a.idAtom.Add(1)
	return id, a.client.Set(ctx, a.genKey(id), value, valid).Err()
}

func (a TokenStore) NewStorePoint(id uint64) TokenStorePoint {
	return TokenStorePoint{
		s:   a,
		key: a.genKey(id),
	}
}

type TokenStorePoint struct {
	s   TokenStore
	key string
}

func (a TokenStorePoint) parsePoint(data []byte, claims interface{}) (bool, error) {
	var parsedData StorePointData
	parsedData.Claims = claims
	return parsedData.Iat == a.s.iat, json.Unmarshal(data, &parsedData)
}

func (a TokenStorePoint) VerifyAndDestroy(ctx context.Context, claims interface{}) (bool, error) {
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return false, err
	}

	return a.parsePoint(value, claims)
}

func (a TokenStorePoint) Verify(ctx context.Context, claims interface{}) (bool, error) {
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return false, err
	}

	return a.parsePoint(value, claims)
}

func (a TokenStorePoint) Destroy(ctx context.Context) error {
	return a.s.client.Del(ctx, a.key).Err()
}
