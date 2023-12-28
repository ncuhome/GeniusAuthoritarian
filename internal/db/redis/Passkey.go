package redis

import (
	"context"
	"encoding/json"
	"time"
)

type PasskeyNamespace string

const (
	PasskeyUser         PasskeyNamespace = "u"
	PasskeyUserRegister PasskeyNamespace = "ur"
	PasskeyLogin        PasskeyNamespace = "l"
)

func NewPasskey(ip string, namespace PasskeyNamespace, identity string) Passkey {
	return Passkey{
		key: keyPasskey.String() + "ip" + ip + string(namespace) + identity,
	}
}

type Passkey struct {
	key string
}

func (p Passkey) StoreSession(ctx context.Context, session any, expire time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return Client.Set(ctx, p.key, data, expire).Err()
}

// ReadSession 读取后自动销毁
func (p Passkey) ReadSession(ctx context.Context, session any) error {
	data, err := Client.GetDel(ctx, p.key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, session)
}
