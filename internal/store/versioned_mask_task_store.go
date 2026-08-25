package store

import "sync"

type VersionedMaskTask struct {
	Status  string
	Version int
}

type VersionedMaskTaskStore struct {
	mu      sync.Mutex
	task    VersionedMaskTask
	effects int
	seen    map[string]bool
}

// NewVersionedMaskTaskStore 构造一个带版本号与效果去重的任务存储。
func NewVersionedMaskTaskStore() *VersionedMaskTaskStore {
	return &VersionedMaskTaskStore{seen: make(map[string]bool)}
}

// Save 以乐观并发方式写入任务状态。
// 仅当传入版本号不早于当前版本号时才接受写入，迟到且版本号更旧的回调会被丢弃，
// 从而保证状态单调推进、重试成功后不会倒退。
// 返回 true 表示写入被接受，false 表示因版本号过期而被忽略。
func (s *VersionedMaskTaskStore) Save(task VersionedMaskTask) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.Version < s.task.Version {
		return false
	}
	s.task = task
	return true
}

// RecordEffect 记录一次对外脱敏动作，按 key 幂等去重。
// 同一个 key 只有首次记录会真正计数，保证多次尝试只产生一次对外效果。
// 返回 true 表示首次记录，false 表示重复记录被忽略。
func (s *VersionedMaskTaskStore) RecordEffect(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen == nil {
		s.seen = make(map[string]bool)
	}
	if s.seen[key] {
		return false
	}
	s.seen[key] = true
	s.effects++
	return true
}

// State 返回当前任务状态与已产生的对外效果次数（按 key 去重后）。
func (s *VersionedMaskTaskStore) State() (VersionedMaskTask, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.task, s.effects
}
