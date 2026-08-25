package store

import "sync"

type MaskedBatchCache struct {
	mu      sync.RWMutex
	batches map[string][]byte
}

func NewMaskedBatchCache() *MaskedBatchCache {
	return &MaskedBatchCache{batches: make(map[string][]byte)}
}

func (c *MaskedBatchCache) Save(id string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batches[id] = payload
}

func (c *MaskedBatchCache) Load(id string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.batches[id]
}
