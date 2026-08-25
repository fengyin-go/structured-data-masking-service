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
}

func (s *VersionedMaskTaskStore) Save(task VersionedMaskTask) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.task = task
}

func (s *VersionedMaskTaskStore) RecordEffect(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.effects++
}

func (s *VersionedMaskTaskStore) State() (VersionedMaskTask, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.task, s.effects
}
