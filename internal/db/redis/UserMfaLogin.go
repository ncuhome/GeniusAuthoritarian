package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var UserMfaLogin = UserMfaLoginHelper{
	key: keyUserMfaLogin.String(),
}

type UserMfaLoginHelper struct {
	key string
}

func (a UserMfaLoginHelper) genKey(uid uint, token string) string {
	return a.key + fmt.Sprint(uid) + "-" + token[:5]
}

func (a UserMfaLoginHelper) Set(uid uint, token string, valid time.Duration, claims interface{}) error {
	v, e := json.Marshal(claims)
	if e != nil {
		return e
	}
	return Client.Set(context.Background(), a.genKey(uid, token), string(v), valid).Err()
}

func (a UserMfaLoginHelper) Get(uid uint, token string, claims interface{}) error {
	value, e := Client.Get(context.Background(), a.genKey(uid, token)).Result()
	if e != nil {
		return e
	}
	return json.Unmarshal([]byte(value), claims)
}
