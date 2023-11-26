package redis

import (
	"context"
	"fmt"
	"time"
)

func NewAccessJwt(uid uint) AccessJwt {
	return AccessJwt{
		key: keyAccessJwt.String() + fmt.Sprint(uid),
	}
}

type AccessJwt struct {
	key string
}

func (a AccessJwt) Set(iat time.Time, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, fmt.Sprint(iat.Unix()), valid).Err()
}

func (a AccessJwt) Pair(iat time.Time) (bool, error) {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return value == fmt.Sprint(iat.Unix()), nil
}

func (a AccessJwt) Clear() error {
	return Client.Del(context.Background(), a.key).Err()
}
