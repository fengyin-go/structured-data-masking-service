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

// Run 在 ctx 取消前最多重试 5 次下游调用。
//
// 关键点：每次调用都传入 ctx，并在每轮开始前检查 ctx.Err()。
// 这样 ctx 超时或取消后：
//   - 不会再发起新的下游调用；
//   - 正在进行的调用会因 Call 内部监听 ctx.Done() 而尽快返回。
//
// 由于 Call 自身已尊重 ctx，Run 不再需要额外的 goroutine + select
// 来“假装”可取消——那种写法会在 ctx 取消后泄漏 goroutine，使其
// 继续跑完整轮重试，正是要修复的问题。
func (r *RetryMaskRunner) Run(ctx context.Context) error {
	for attempt := 0; attempt < 5; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := r.client.Call(ctx); err != nil {
			return err
		}
	}
	return nil
}
