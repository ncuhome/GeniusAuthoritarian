package redis

import (
	"context"
	"fmt"
	"time"
)

var MfaEnable = MfaEnableHelper{
	key: keyMfaEnable.String(),
}

type MfaEnableHelper struct {
	key string
}

func (a MfaEnableHelper) Set(uid uint, secret string, valid time.Duration) error {
	return Client.Set(context.Background(), a.key+fmt.Sprint(uid), secret, valid).Err()
}

func (a MfaEnableHelper) Get(uid uint) (string, error) {
	return Client.Get(context.Background(), a.key+fmt.Sprint(uid)).Result()
}

func (a MfaEnableHelper) Del(uid uint) error {
	return Client.Del(context.Background(), a.key+fmt.Sprint(uid)).Err()
}
