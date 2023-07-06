package backoff

import (
	"github.com/Mmx233/tool"
	"sync"
	"time"
)

type Conf struct {
	Content       func() error
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
		lock:     &sync.Mutex{},
	}
}

// Backoff 错误重试 积分退避算法
type Backoff struct {
	f        func() error
	retry    time.Duration
	maxRetry time.Duration

	lock    *sync.Mutex
	running bool
}

func (a Backoff) Start() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.running {
		a.running = true
		go a.Worker()
	}
}

// Worker
// 请注意,此处使用的是普通接收器,当 worker 重新运行时 retry 会被重置
func (a Backoff) Worker() {
	defer tool.Recover()

	for {
		e := a.f()
		if e == nil {
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

	a.running = false
}
