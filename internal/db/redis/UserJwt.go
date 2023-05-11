package redis

import (
	"context"
	"fmt"
	"time"
)

var UserJwt = UserJwtHelper{
	key:         keyUserJwt.String(),
	storeLength: 5,
}

type UserJwtHelper struct {
	key         string
	storeLength int
}

func (a UserJwtHelper) userKey(uid uint) string {
	return a.key + fmt.Sprint(uid)
}

func (a UserJwtHelper) Set(uid uint, token string, valid time.Duration) error {
	return Client.Set(context.Background(), a.userKey(uid), token[:a.storeLength], valid).Err()
}

func (a UserJwtHelper) Pair(uid uint, token string) (bool, error) {
	value, e := Client.Get(context.Background(), a.userKey(uid)).Result()
	if e != nil {
		if e == Nil {
			e = nil
		}
		return false, e
	}
	return value == token[:a.storeLength], nil
}
