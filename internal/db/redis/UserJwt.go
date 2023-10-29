package redis

import (
	"context"
	"fmt"
	"time"
)

func NewUserJwt(uid uint) UserJwt {
	return UserJwt{
		key:           keyUserJwt.String() + fmt.Sprint(uid),
		compareLength: 5,
	}
}

type UserJwt struct {
	key           string
	compareLength int
}

func (a UserJwt) Set(token string, valid time.Duration) error {
	return Client.Set(context.Background(), a.key, token[:a.compareLength], valid).Err()
}

func (a UserJwt) Pair(token string) (bool, error) {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}
	return value == token[:a.compareLength], nil
}

func (a UserJwt) Clear() error {
	return Client.Del(context.Background(), a.key).Err()
}
