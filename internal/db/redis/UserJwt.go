package redis

import (
	"context"
	"fmt"
	"time"
)

var UserJwt = UserJwtHelper{
	key:           keyUserJwt.String(),
	compareLength: 5,
}

type UserJwtHelper struct {
	key           string
	compareLength int
}

func (a UserJwtHelper) userKey(uid uint) string {
	return a.key + fmt.Sprint(uid)
}

func (a UserJwtHelper) Set(uid uint, token string, valid time.Duration) error {
	return Client.Set(context.Background(), a.userKey(uid), token[:a.compareLength], valid).Err()
}

func (a UserJwtHelper) Pair(uid uint, token string) (bool, error) {
	value, e := Client.Get(context.Background(), a.userKey(uid)).Result()
	if e != nil {
		if e == Nil {
			e = nil
		}
		return false, e
	}
	return value == token[:a.compareLength], nil
}

func (a UserJwtHelper) Clear(uid uint) error {
	return Client.Del(context.Background(), a.userKey(uid)).Err()
}
