package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
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
	keySshDevSub
	keyPasskey
	keyU2F
)

var (
	idThirdPartyLogin,
	idMfaLogin,
	idU2F atomic.Uint64
)
