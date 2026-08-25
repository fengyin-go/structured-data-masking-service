package service_test

import (
	"testing"

	"datamasking/internal/service"
	"datamasking/internal/store"
)

func TestRuleSummaryUsesRequestSnapshot(t *testing.T) {
	st := store.NewRuleSnapshotStore([]string{"email", "phone"})
	summary := service.NewRuleSnapshotSummary(st)
	started := make(chan struct{})
	release := make(chan struct{})
	result := summary.SummarizeAfter(started, release)
	<-started
	st.Replace(0, "identity")
	close(release)
	if got := <-result; got != "email,phone" {
		t.Errorf("in-flight rule summary changed with a later update: got %q", got)
	}
}
