package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func NewPasskey() Passkey {
	return Passkey{}
}

type Passkey struct {
}

func (p Passkey) key(ip string) string {
	return keyPasskey.String() + "ip" + ip
}

func (p Passkey) store(ctx context.Context, key string, session any, expire time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, expire).Err()
}

func (p Passkey) StoreSession(ctx context.Context, ip string, session any, expire time.Duration) error {
	return p.store(ctx, p.key(ip), session, expire)
}

func (p Passkey) read(ctx context.Context, key string, session any) error {
	data, err := Client.GetDel(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, session)
}

// ReadSession 读取后自动销毁
func (p Passkey) ReadSession(ctx context.Context, ip string, session any) error {
	return p.read(ctx, p.key(ip), session)
}

func (p Passkey) NewUser(id uint) UserPasskey {
	return UserPasskey{
		p:   p,
		key: keyPasskey.String() + "u" + fmt.Sprint(id),
	}
}

type UserPasskey struct {
	p   Passkey
	key string
}

func (u UserPasskey) StoreSession(ctx context.Context, session any, expire time.Duration) error {
	return u.p.store(ctx, u.key, session, expire)
}

func (u UserPasskey) ReadSession(ctx context.Context, session any) error {
	return u.p.read(ctx, u.key, session)
}
