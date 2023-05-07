package redis

import (
	"context"
)

var AppCode = AppCodeHelper{
	key: keyAppCode.String(),
}

type AppCodeHelper struct {
	key string
}

func (a AppCodeHelper) Add(data ...string) error {
	return Client.SAdd(context.Background(), a.key, data).Err()
}

func (a AppCodeHelper) Load() ([]string, error) {
	var t []string
	cmd := Client.SMembers(context.Background(), a.key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	if e := cmd.ScanSlice(&t); e != nil {
		return nil, e
	}
	if len(t) == 0 {
		return nil, Nil
	}
	return t, nil
}
