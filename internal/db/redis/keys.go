package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type keys uint8

const Nil = redis.Nil

func (a keys) String() string {
	return fmt.Sprint(uint8(a)) + "-"
}

const (
	keyUserIdentityCode keys = iota
	keyThirdPartyLogin
	keyAppCode
	keyUserJwt
	keyUserMfaLogin
	keyMfaEnable
	keySms
	keySyncStat
)
