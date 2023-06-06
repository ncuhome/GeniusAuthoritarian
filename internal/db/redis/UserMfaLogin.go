package redis

import (
	"context"
	"fmt"
	"time"
)

var UserMfaLogin = UserMfaLoginHelper{
	key:           keyUserMfaLogin.String(),
	compareLength: 5,
}

type UserMfaLoginHelper struct {
	key           string
	compareLength int
}

func (a UserMfaLoginHelper) Set(uid uint, token string, valid time.Duration) error {
	return Client.Set(context.Background(), a.key+fmt.Sprint(uid), token[:a.compareLength], valid).Err()
}

func (a UserMfaLoginHelper) Verify(uid uint, token string) (bool, error) {
	value, e := Client.Get(context.Background(), a.key+fmt.Sprint(uid)).Result()
	if e != nil {
		if e == Nil {
			e = nil
		}
		return false, e
	}
	return value == token[:a.compareLength], nil
}
