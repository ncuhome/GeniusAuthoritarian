package redis

import (
	"context"
)

var SiteWhiteList = SiteWhiteListHelper{
	key: keySiteWhiteList.String(),
}

type SiteWhiteListHelper struct {
	key string
}

func (a SiteWhiteListHelper) Add(data ...string) error {
	return Client.SAdd(context.Background(), a.key, data).Err()
}

func (a SiteWhiteListHelper) Load() ([]string, error) {
	var t []string
	cmd := Client.SMembers(context.Background(), a.key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return t, cmd.ScanSlice(&t)
}
