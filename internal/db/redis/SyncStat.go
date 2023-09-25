package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewSyncStat(name string) SyncStat {
	key := keySyncStat.String() + name
	return SyncStat{
		key:     key,
		lockKey: key + "-lock",
	}
}

type SyncStat struct {
	key     string
	lockKey string

	lockMark string
}

func (a SyncStat) TryLock(ctx context.Context, expire time.Duration) (bool, error) {
	a.lockMark = fmt.Sprint(time.Now().UnixNano())
	return Client.SetNX(ctx, a.lockKey, a.lockMark, expire).Result()
}

func (a SyncStat) MustLock(ctx context.Context, expire time.Duration) error {
	var count uint8
	for ; count < 255; count++ {
		ok, err := a.TryLock(ctx, expire)
		if err != nil {
			return err
		} else if ok {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}

	return errors.New("wait for sync lock timeout")
}

func (a SyncStat) Unlock(ctx context.Context) error {
	mark, err := Client.Get(ctx, a.lockKey).Result()
	if err != nil {
		return err
	} else if mark != a.lockMark {
		return nil
	}
	return Client.Del(ctx, a.lockKey).Err()
}

func (a SyncStat) SetSuccess(ctx context.Context, expire time.Duration) error {
	return Client.Set(ctx, a.key, "1", expire).Err()
}

func (a SyncStat) Succeed(ctx context.Context) (bool, error) {
	err := Client.Get(ctx, a.key).Err()
	if err == redis.Nil {
		return false, nil
	}
	return true, err
}
