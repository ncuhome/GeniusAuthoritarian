package singleton

import (
	"context"
	"sync/atomic"
)

type Instance interface {
	Run(ctx context.Context) error
}

func New(fn func(ctx context.Context) error) Instance {
	return &Singleton{Fn: fn}
}

type Singleton struct {
	Fn      func(ctx context.Context) error
	Running atomic.Bool
}

func (s *Singleton) Run(ctx context.Context) error {
	if s.Running.CompareAndSwap(false, true) {
		defer s.Running.Store(false)
		return s.Fn(ctx)
	}
	return ErrAlreadyRunning
}
