package service

import (
	"strings"

	"datamasking/internal/store"
)

type RuleSnapshotSummary struct {
	store *store.RuleSnapshotStore
}

func NewRuleSnapshotSummary(st *store.RuleSnapshotStore) *RuleSnapshotSummary {
	return &RuleSnapshotSummary{store: st}
}

func (s *RuleSnapshotSummary) SummarizeAfter(started chan<- struct{}, release <-chan struct{}) <-chan string {
	rules := s.store.Snapshot()
	result := make(chan string, 1)
	go func() {
		close(started)
		<-release
		result <- strings.Join(rules, ",")
		close(result)
	}()
	return result
}
