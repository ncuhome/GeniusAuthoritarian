package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
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

func (a *SyncStat) TryLock(ctx context.Context, expire time.Duration) (bool, error) {
	a.lockMark = fmt.Sprint(time.Now().UnixNano())
	return Client.SetNX(ctx, a.lockKey, a.lockMark, expire).Result()
}

func (a *SyncStat) MustLock(ctx context.Context, expire time.Duration) error {
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

func (a *SyncStat) Unlock(ctx context.Context) error {
	mark, err := Client.Get(ctx, a.lockKey).Result()
	if err != nil {
		return err
	} else if mark != a.lockMark {
		return nil
	}
	return Client.Del(ctx, a.lockKey).Err()
}

func (a *SyncStat) SetSuccess(ctx context.Context, expire time.Duration) error {
	return Client.Set(ctx, a.key, "1", expire).Err()
}

func (a *SyncStat) Succeed(ctx context.Context) (bool, error) {
	err := Client.Get(ctx, a.key).Err()
	if err == redis.Nil {
		return false, nil
	}
	return true, err
}

// Inject 注入 backoff 内容函数，使其支持分布式锁与成功跳过
func (a *SyncStat) Inject(schedule cron.Schedule, f func() error) func() error {
	return func() error {
		ok, err := a.Succeed(context.Background())
		if err != nil {
			return err
		} else if ok {
			return nil
		}

		if err = a.MustLock(context.Background(), time.Second*120); err != nil {
			return err
		}
		defer a.Unlock(context.Background())

		if err = f(); err != nil {
			return err
		} else {
			next := schedule.Next(time.Now())
			expire := next.Sub(time.Now()) - time.Second*5
			if expire > 0 {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_ = a.SetSuccess(ctx, expire)
			}
		}

		return nil
	}
}
