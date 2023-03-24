package redis

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var Jwt = JwtHelper{
	key: keyJwt.String(),
}

type JwtHelper struct {
	key string
	id  atomic.Uint64
}

func (a *JwtHelper) authPointKey(id uint64) string {
	return a.key + "ap-" + fmt.Sprint(id)
}

func (a *JwtHelper) NewAuthPoint(unix int64, valid time.Duration) (id uint64, e error) {
	id = a.id.Add(1)
	e = Client.Set(context.Background(), a.authPointKey(id), fmt.Sprint(unix), valid).Err()
	return
}

func (a *JwtHelper) VerifyAuthPoint(id uint64, unix int64) (bool, error) {
	v, e := Client.Get(context.Background(), a.authPointKey(id)).Result()
	if e != nil {
		if e == Nil {
			e = nil
		}
		return false, e
	}
	return v == fmt.Sprint(unix), nil
}
