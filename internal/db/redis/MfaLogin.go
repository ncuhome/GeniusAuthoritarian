package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var MfaLogin = MfaLoginHelper{
	key: keyUserMfaLogin.String(),
}

type MfaLoginHelper struct {
	key string
}

func (a MfaLoginHelper) genKey(uid uint, token string) string {
	return a.key + fmt.Sprint(uid) + "-" + token[:5]
}

func (a MfaLoginHelper) Set(uid uint, token string, valid time.Duration, claims interface{}) error {
	v, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return Client.Set(context.Background(), a.genKey(uid, token), string(v), valid).Err()
}

func (a MfaLoginHelper) Get(uid uint, token string, claims interface{}) error {
	value, err := Client.Get(context.Background(), a.genKey(uid, token)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), claims)
}
