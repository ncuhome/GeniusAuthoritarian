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

func (a *ThirdPartyLoginHelper) NewLoginPoint(unix int64, valid time.Duration, claims interface{}) (id uint64, err error) {
	claimsRaw, err := json.Marshal(claims)
	if err != nil {
		return 0, err
	}

	loginPoint := LoginPoint{
		Unix: unix,
		Data: claimsRaw,
	}
	value, err := json.Marshal(loginPoint)
	if err != nil {
		return 0, err
	}

	id = a.id.Add(1)
	err = Client.Set(context.Background(), a.loginPointKey(id), value, valid).Err()
	return
}

func (a *ThirdPartyLoginHelper) VerifyLoginPoint(id uint64, unix int64, claims interface{}) (bool, error) {
	value, err := Client.GetDel(context.Background(), a.loginPointKey(id)).Result()
	if err != nil {
		if err == Nil {
			err = nil
		}
		return false, err
	}

	var loginPoint LoginPoint
	if err = json.Unmarshal([]byte(value), &loginPoint); err != nil {
		return false, err
	}
	return loginPoint.Unix == unix, json.Unmarshal(loginPoint.Data, claims)
}

func (a *ThirdPartyLoginHelper) DelLoginPoint(id uint64) error {
	return Client.Del(context.Background(), a.loginPointKey(id)).Err()
}
