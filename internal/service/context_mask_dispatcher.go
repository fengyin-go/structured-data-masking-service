package service

import (
	"context"

	"datamasking/internal/store"
)

type ContextMaskDispatcher struct {
	sink *store.DelayedMaskSink
}

func NewContextMaskDispatcher(sink *store.DelayedMaskSink) *ContextMaskDispatcher {
	return &ContextMaskDispatcher{sink: sink}
}

func (d *ContextMaskDispatcher) Dispatch(ctx context.Context, value string) error {
	return d.sink.Write(context.Background(), value)
}
