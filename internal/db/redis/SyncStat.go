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

func (a *SyncStat) ShouldLock(ctx context.Context, expire time.Duration) (bool, error) {
	var count uint16
	for ; count < 500; count++ {
		ok, err := a.Succeed(ctx)
		if err != nil {
			return false, err
		} else if ok {
			return false, nil
		}

		ok, err = a.TryLock(ctx, expire)
		if err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
		time.Sleep(time.Millisecond * 35)
	}

	return false, errors.New("wait for sync lock timeout")
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
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return true, err
}

// Inject 注入 backoff 内容函数，使其支持分布式锁与成功跳过
func (a *SyncStat) Inject(schedule cron.Schedule, f func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ok, err := a.Succeed(ctx)
		if err != nil {
			return err
		} else if ok {
			return nil
		}

		locked, err := a.ShouldLock(context.Background(), time.Second*120)
		if err != nil {
			return err
		} else if !locked {
			return nil
		}

		defer a.Unlock(ctx)

		if err = f(ctx); err != nil {
			return err
		} else {
			next := schedule.Next(time.Now())
			expire := next.Sub(time.Now()) - time.Second*5
			if expire > 0 {
				ctx, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()
				_ = a.SetSuccess(ctx, expire)
			}
		}

		return nil
	}
}
