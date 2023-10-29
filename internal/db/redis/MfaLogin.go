package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func NewMfaLogin(uid uint, token string) MfaLogin {
	return MfaLogin{
		key: keyUserMfaLogin.String() + fmt.Sprint(uid) + "-" + token[:5],
	}
}

type MfaLogin struct {
	key string
}

func (a MfaLogin) Set(valid time.Duration, claims interface{}) error {
	v, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return Client.Set(context.Background(), a.key, string(v), valid).Err()
}

func (a MfaLogin) Get(claims interface{}) error {
	value, err := Client.Get(context.Background(), a.key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), claims)
}
