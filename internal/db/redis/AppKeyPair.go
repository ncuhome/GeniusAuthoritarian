package redis

import (
	"context"
	"time"
)

func NewAppKeyPair(appCode string) AppKeyPair {
	return AppKeyPair{
		key: keyAppKeyPair.String() + appCode,
	}
}

type AppKeyPair struct {
	key string
}

func (a AppKeyPair) Cache(ctx context.Context, appSecret string) error {
	return Client.Set(ctx, a.key, appSecret, time.Hour).Err()
}

func (a AppKeyPair) Read(ctx context.Context) (string, error) {
	return Client.GetEx(ctx, a.key, time.Hour).Result()
}
