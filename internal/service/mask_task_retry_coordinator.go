package service

import (
	"datamasking/internal/store"
)

// MaskTaskRetryCoordinator 编排脱敏任务的首次执行与重试，处理首次回调延迟到达的场景。
//
// 关注的不变式：
//  1. 重试成功后状态不可倒退——迟到且版本号更旧的回调不得覆盖更新的成功结果。
//  2. 同一对外脱敏动作多次尝试只产生一次效果——按稳定的 key 幂等去重。
type MaskTaskRetryCoordinator struct {
	store *store.VersionedMaskTaskStore
}

func NewMaskTaskRetryCoordinator(st *store.VersionedMaskTaskStore) *MaskTaskRetryCoordinator {
	return &MaskTaskRetryCoordinator{store: st}
}

// RunWithDelayedCallback 演绎“首次执行延迟、触发重试、旧回调晚到”的场景。
//
// 时序：
//  - 尝试 1：状态置 running(v1)，记录对外动作 mask-attempt。
//  - 其回调被挂起（等待 release）。
//  - 尝试 2（重试）：状态置 done(v2)，记录对外动作 mask-attempt。
//  - 调用方稍后释放 release，模拟迟到的首次回调。
//
// 由于 store.Save 按版本号单调推进，迟到的 v1 回调写入会被丢弃，状态停留在 done(v2)；
// 对外动作的 key 在两次尝试中相同，RecordEffect 仅生效一次。
func (c *MaskTaskRetryCoordinator) RunWithDelayedCallback() (chan<- struct{}, <-chan struct{}) {
	release := make(chan struct{})
	done := make(chan struct{})

	// 首次执行：状态置 running，版本号 1。
	c.store.Save(store.VersionedMaskTask{Status: "running", Version: 1})
	// 对外脱敏动作，跨尝试复用同一 key，保证只生效一次。
	c.store.RecordEffect("mask-attempt")

	// 首次执行的回调被延迟（模拟网络/异步回调晚到）。
	go func() {
		<-release
		// 迟到的回调仍以版本号 1 写回 running。
		// 此时已完成的重试结果为 done(v2)，store.Save 拒绝版本号更旧的写入，状态不倒退。
		c.store.Save(store.VersionedMaskTask{Status: "running", Version: 1})
		c.store.RecordEffect("mask-attempt")
		close(done)
	}()

	// 重试成功：状态置 done，版本号 2，严格高于首次执行，覆盖 running(v1)。
	c.store.Save(store.VersionedMaskTask{Status: "done", Version: 2})
	// 同一对外动作，去重后不重复计数。
	c.store.RecordEffect("mask-attempt")

	return release, done
}
