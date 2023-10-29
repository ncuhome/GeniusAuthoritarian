package tokenStore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
	"time"
)

func NewTokenStore(Client *redis.Client, idAtom *atomic.Uint64, keyPrefix string) TokenStore {
	return TokenStore{
		client:    Client,
		keyPrefix: keyPrefix,
		idAtom:    idAtom,
	}
}

type TokenStore struct {
	client    *redis.Client
	keyPrefix string
	idAtom    *atomic.Uint64
}

func (a TokenStore) genKey(id uint64) string {
	return a.keyPrefix + fmt.Sprint(id)
}

func (a TokenStore) CreateStorePoint(ctx context.Context, iat time.Time, valid time.Duration, claims interface{}) (uint64, error) {
	var value []byte
	if claims == nil {
		value = []byte{'1'}
	} else {
		var err error
		value, err = json.Marshal(&StorePointData{
			Iat:    iat.Unix(),
			Claims: claims,
		})
		if err != nil {
			return 0, err
		}
	}

	id := a.idAtom.Add(1)
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
	value, err := a.s.client.GetDel(ctx, a.key).Bytes()
	if err != nil {
		return err
	}

	return a.parsePoint(iat, value, claims)
}

func (a Point) Destroy(ctx context.Context) error {
	return a.s.client.Del(ctx, a.key).Err()
}
