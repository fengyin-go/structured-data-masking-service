package service

import (
	"context"

	"datamasking/internal/store"
)

type RetryMaskRunner struct {
	client *store.RetryMaskClient
}

func NewRetryMaskRunner(client *store.RetryMaskClient) *RetryMaskRunner {
	return &RetryMaskRunner{client: client}
}

func (r *RetryMaskRunner) Run(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		for attempt := 0; attempt < 5; attempt++ {
			if err := r.client.Call(context.Background()); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
