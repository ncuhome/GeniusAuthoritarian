package redis

import (
	"context"
	"fmt"
	"time"
)

// NewUserToken 用户后台 token
func NewUserToken(uid uint) UserToken {
	return UserToken{
		key: keyUserToken.String() + fmt.Sprint(uid),
	}
}

type UserToken struct {
	key string
}

func (a UserToken) Set(iat time.Time, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, fmt.Sprint(iat.Unix()), valid).Err()
}

func (a UserToken) Pair(iat time.Time) (bool, error) {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return value == fmt.Sprint(iat.Unix()), nil
}

func (a UserToken) Clear() error {
	return Client.Del(context.Background(), a.key).Err()
}
