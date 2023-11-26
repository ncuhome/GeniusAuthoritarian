package redis

import (
	"context"
	"fmt"
	"time"
)

func NewRefreshJwt(uid uint) RefreshJwt {
	return RefreshJwt{
		key: keyRefreshJwt.String() + fmt.Sprint(uid),
	}
}

type RefreshJwt struct {
	key string
}

func (a RefreshJwt) Set(iat time.Time, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, fmt.Sprint(iat), (valid)*time.Second).Err()
}

func (a RefreshJwt) Pair(iat time.Time) (bool, error) {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return value == fmt.Sprint(iat), nil
}

func (a RefreshJwt) Clear() error {
	return Client.Del(context.Background(), a.key).Err()
}
