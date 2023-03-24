package feishu

import (
	"sync"
	"time"
)

func newTenant() *tenantTokenCache {
	return &tenantTokenCache{}
}

type tenantTokenCache struct {
	sync.RWMutex
	Token    string
	ExpireAt int64
}

func (t *tenantTokenCache) Load() (string, bool) {
	t.RLock()
	if t.Token != "" {
		if t.ExpireAt-30 > time.Now().Unix() {
			defer t.RUnlock()
			return t.Token, true
		} else {
			t.RUnlock()
			t.Lock()
			defer t.Unlock()
			t.Token = ""
			return "", false
		}
	}
	t.RUnlock()
	return "", false
}

func (t *tenantTokenCache) Set(token string, expire int64) {
	t.Lock()
	defer t.Unlock()
	t.Token = token
	t.ExpireAt = time.Now().Add(time.Second * time.Duration(expire)).Unix()
}
