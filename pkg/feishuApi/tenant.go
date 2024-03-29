package feishuApi

import (
	"sync"
	"time"
)

func NewTenant(fs *Fs) *TenantToken {
	return &TenantToken{
		fs: fs,
	}
}

type TenantToken struct {
	fs *Fs

	sync.Mutex
	Token    string
	ExpireAt int64
}

func (t *TenantToken) Load() (string, error) {
	t.Lock()
	defer t.Unlock()

	if t.Token != "" && t.ExpireAt-30 > time.Now().Unix() {
		return t.Token, nil
	}

	t.Token = ""
	tokenRes, err := t.fs.GetTenantAccessToken()
	if err != nil {
		return "", err
	}

	t.Token = tokenRes.TenantAccessToken
	t.ExpireAt = time.Now().Add(time.Second * time.Duration(tokenRes.Expire)).Unix()
	return t.Token, nil
}
