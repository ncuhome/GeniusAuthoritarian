package redis

import (
	"context"
	"time"
)

var Sms = SmsHelper{
	key: keySms.String(),
}

type SmsHelper struct {
	key string
}

func (a SmsHelper) genKey(phone string) string {
	return a.key + "-" + phone
}

func (a SmsHelper) TryLock(phone string) (bool, error) {
	err := Client.SetNX(context.Background(), phone, a.genKey(phone), time.Minute).Err()
	if err == Nil {
		return false, nil
	}
	return err == nil, err
}

func (a SmsHelper) IsLocked(phone string) (bool, error) {
	err := Client.Get(context.Background(), a.genKey(phone)).Err()
	if err == Nil {
		return false, nil
	}
	return err == nil, err
}

func (a SmsHelper) UnLock(phone string) error {
	return Client.Del(context.Background(), a.genKey(phone)).Err()
}
