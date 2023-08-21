package backoff

import (
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"sync/atomic"
	"time"
)

// Backoff 错误重试 积分退避算法
type Backoff struct {
	f        func() error
	retry    time.Duration
	maxRetry time.Duration

	running *atomic.Bool
}

type Conf struct {
	Content func() error
	// 最大重试等待时间
	MaxRetryDelay time.Duration
}

func New(c Conf) Backoff {
	if c.Content == nil {
		panic("content function required")
	}
	if c.MaxRetryDelay == 0 {
		c.MaxRetryDelay = time.Minute * 20
	}

	return Backoff{
		f:        c.Content,
		retry:    time.Second,
		maxRetry: c.MaxRetryDelay,
		running:  &atomic.Bool{},
	}
}

func (a Backoff) Start() {
	if !a.running.Load() {
		if a.running.CompareAndSwap(false, true) {
			go a.Worker()
		}
	}
}

// Worker
// 请注意,此处使用的是普通接收器,当 worker 重新运行时 retry 会被重置
func (a Backoff) Worker() {
	for {
		errChan := make(chan error)
		go func() {
			defer func() {
				if p := tool.Recover(); p != nil {
					errChan <- errors.New(fmt.Sprint(p))
				}
			}()
			errChan <- a.f()
		}()
		if err := <-errChan; err == nil {
			break
		}

		time.Sleep(a.retry)

		if a.retry < a.maxRetry {
			a.retry = a.retry << 1
			if a.retry > a.maxRetry {
				a.retry = a.maxRetry
			}
		}
	}

	a.running.Store(false)
}
