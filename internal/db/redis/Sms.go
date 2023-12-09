package redis

import (
	"context"
	"errors"
	"time"
)

func NewSms(phone string) Sms {
	return Sms{
		key: keySms.String() + "-" + phone,
	}
}

type Sms struct {
	key string
}

func (a Sms) TryLock() (bool, error) {
	return Client.SetNX(context.Background(), a.key, "1", time.Minute).Result()
}

func (a Sms) IsLocked() (bool, error) {
	err := Client.Get(context.Background(), a.key).Err()
	if errors.Is(err, Nil) {
		return false, nil
	}
	return err == nil, err
}

func (a Sms) UnLock() error {
	return Client.Del(context.Background(), a.key).Err()
}
