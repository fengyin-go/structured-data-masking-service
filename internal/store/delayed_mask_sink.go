package store

import (
	"context"
	"sync"
	"time"
)

type DelayedMaskSink struct {
	delay time.Duration
	mu    sync.Mutex
	ctx   context.Context
}

func NewDelayedMaskSink(delay time.Duration) *DelayedMaskSink {
	return &DelayedMaskSink{delay: delay}
}

func (s *DelayedMaskSink) Write(ctx context.Context, value string) error {
	s.mu.Lock()
	if s.ctx == nil {
		s.ctx = ctx
	}
	active := s.ctx
	s.mu.Unlock()

	select {
	case <-active.Done():
		return active.Err()
	case <-time.After(s.delay):
		return nil
	}
}
