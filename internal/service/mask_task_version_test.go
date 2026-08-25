package service_test

import (
	"testing"

	"datamasking/internal/service"
	"datamasking/internal/store"
)

func TestMaskTaskRetryRejectsStaleCompletion(t *testing.T) {
	st := &store.VersionedMaskTaskStore{}
	coordinator := service.NewMaskTaskRetryCoordinator(st)
	release, callbackDone := coordinator.RunWithDelayedCallback()
	close(release)
	<-callbackDone
	task, effects := st.State()
	if task.Status != "done" || task.Version != 2 {
		t.Errorf("late first attempt moved a completed mask task backward: %+v", task)
	}
	if effects != 1 {
		t.Errorf("mask retry repeated its external effect %d times", effects)
	}
}
