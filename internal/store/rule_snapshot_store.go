package store

import "sync"

type RuleSnapshotStore struct {
	mu    sync.RWMutex
	rules []string
}

func NewRuleSnapshotStore(rules []string) *RuleSnapshotStore {
	return &RuleSnapshotStore{rules: rules}
}

func (s *RuleSnapshotStore) Snapshot() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rules
}

func (s *RuleSnapshotStore) Replace(index int, rule string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[index] = rule
}
