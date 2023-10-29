package redis

import (
	"context"
	"fmt"
	"time"
)

func NewMfaEnable(uid uint) MfaEnable {
	return MfaEnable{
		key: keyMfaEnable.String() + fmt.Sprint(uid),
	}
}

type MfaEnable struct {
	key string
}

func (a MfaEnable) Set(secret string, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, secret, valid).Err()
}

func (a MfaEnable) Get() (string, error) {
	return Client.Get(context.Background(), a.key).Result()
}

func (a MfaEnable) Del() error {
	return Client.Del(context.Background(), a.key).Err()
}
