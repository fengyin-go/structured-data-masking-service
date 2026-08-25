package store

import "sync"

type RuleSnapshotStore struct {
	mu    sync.RWMutex
	rules []string
}

func NewRuleSnapshotStore(rules []string) *RuleSnapshotStore {
	return &RuleSnapshotStore{rules: rules}
}

// Snapshot 返回当前规则的一个独立副本。
// 必须拷贝底层数组：调用方持有快照后，Replace 仍会就地修改 store 内部
// 的 rules。若直接返回 s.rules，处理中的请求在放行后会读到被后续更新
// 改写的内容，造成串批次。拷贝后快照固定为请求开始时的内容，后续
// Replace 只影响下一次 Snapshot。
func (s *RuleSnapshotStore) Snapshot() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.rules))
	copy(out, s.rules)
	return out
}

func (s *RuleSnapshotStore) Replace(index int, rule string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[index] = rule
}
