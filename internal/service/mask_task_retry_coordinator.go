package service

import (
	"fmt"

	"datamasking/internal/store"
)

type MaskTaskRetryCoordinator struct {
	store *store.VersionedMaskTaskStore
}

func NewMaskTaskRetryCoordinator(st *store.VersionedMaskTaskStore) *MaskTaskRetryCoordinator {
	return &MaskTaskRetryCoordinator{store: st}
}

func (c *MaskTaskRetryCoordinator) RunWithDelayedCallback() (chan<- struct{}, <-chan struct{}) {
	release := make(chan struct{})
	done := make(chan struct{})
	c.store.Save(store.VersionedMaskTask{Status: "running", Version: 1})
	c.store.RecordEffect(fmt.Sprintf("mask-attempt-%d", 1))
	go func() {
		<-release
		c.store.Save(store.VersionedMaskTask{Status: "running", Version: 1})
		c.store.RecordEffect("mask-delayed-callback")
		close(done)
	}()
	c.store.Save(store.VersionedMaskTask{Status: "done", Version: 2})
	c.store.RecordEffect(fmt.Sprintf("mask-attempt-%d", 2))
	return release, done
}
