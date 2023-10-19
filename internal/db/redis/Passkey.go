package redis

import (
	"context"
	"encoding/json"
	"time"
)

func NewPasskey() Passkey {
	return Passkey{}
}

type Passkey struct {
}

func (p Passkey) key(ip string) string {
	return keyPasskey.String() + "ip-" + ip
}

func (p Passkey) StoreSession(ctx context.Context, ip string, session any, expire time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return Client.Set(ctx, p.key(ip), data, expire).Err()
}

// ReadSession 读取后自动销毁
func (p Passkey) ReadSession(ctx context.Context, ip string, session any) error {
	data, err := Client.GetDel(ctx, p.key(ip)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, session)
}
