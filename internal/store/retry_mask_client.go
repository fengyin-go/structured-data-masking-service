package store

import (
	"context"
	"sync"
	"time"
)

type RetryMaskClient struct {
	mu    sync.Mutex
	delay time.Duration
	calls int
}

func NewRetryMaskClient(delay time.Duration) *RetryMaskClient {
	return &RetryMaskClient{delay: delay}
}

func (c *RetryMaskClient) Call(ctx context.Context) error {
	c.mu.Lock()
	c.calls++
	c.mu.Unlock()
	<-time.After(c.delay)
	return nil
}

func (c *RetryMaskClient) Calls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls
}
