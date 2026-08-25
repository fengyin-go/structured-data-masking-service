package store

import (
	"sync"

	"datamasking/internal/model"
)

type MemoryStore struct {
	mu          sync.RWMutex
	maskRules   map[string]*model.MaskRule
	dataSources map[string]*model.DataSource
	maskTasks   map[string]*model.MaskTask
	maskRecords map[string]*model.MaskRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		maskRules:   make(map[string]*model.MaskRule),
		dataSources: make(map[string]*model.DataSource),
		maskTasks:   make(map[string]*model.MaskTask),
		maskRecords: make(map[string]*model.MaskRecord),
	}
}

var _ Store = (*MemoryStore)(nil)
