package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

var ThirdPartyLogin = ThirdPartyLoginHelper{
	key: keyThirdPartyLogin.String(),
}

type ThirdPartyLoginHelper struct {
	key string
	id  atomic.Uint64
}

func (a *ThirdPartyLoginHelper) loginPointKey(id uint64) string {
	return a.key + "ap-" + fmt.Sprint(id)
}

type LoginPoint struct {
	Unix int64 `json:"unix"`
	Data json.RawMessage
}

func (a *ThirdPartyLoginHelper) NewLoginPoint(unix int64, valid time.Duration, claims interface{}) (id uint64, e error) {
	claimsRaw, e := json.Marshal(claims)
	if e != nil {
		return 0, e
	}

	loginPoint := LoginPoint{
		Unix: unix,
		Data: claimsRaw,
	}
	value, e := json.Marshal(loginPoint)
	if e != nil {
		return 0, e
	}

	id = a.id.Add(1)
	e = Client.Set(context.Background(), a.loginPointKey(id), value, valid).Err()
	return
}

func (a *ThirdPartyLoginHelper) VerifyLoginPoint(id uint64, unix int64, claims interface{}) (bool, error) {
	value, e := Client.Get(context.Background(), a.loginPointKey(id)).Result()
	if e != nil {
		if e == Nil {
			e = nil
		}
		return false, e
	}

	var loginPoint LoginPoint
	if e = json.Unmarshal([]byte(value), &loginPoint); e != nil {
		return false, e
	}
	return loginPoint.Unix == unix, json.Unmarshal(loginPoint.Data, claims)
}

func (a *ThirdPartyLoginHelper) DelLoginPoint(id uint64) error {
	return Client.Del(context.Background(), a.loginPointKey(id)).Err()
}
