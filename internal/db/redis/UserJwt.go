package redis

import (
	"context"
	"fmt"
	"time"
)

func NewUserJwt(uid uint) UserJwt {
	return UserJwt{
		key: keyUserJwt.String() + fmt.Sprint(uid),
	}
}

type UserJwt struct {
	key string
}

func (a UserJwt) Set(iat time.Time, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, fmt.Sprint(iat.Unix()), valid).Err()
}

func (a UserJwt) Pair(iat time.Time) (bool, error) {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return value == fmt.Sprint(iat.Unix()), nil
}

func (a UserJwt) Clear() error {
	return Client.Del(context.Background(), a.key).Err()
}
