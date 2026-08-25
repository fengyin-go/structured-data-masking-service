package store

import "sync"

type MaskPlan struct {
	Fields []string
	Ready  bool
}

type MaskPlanCache struct {
	mu    sync.RWMutex
	plans map[string]*MaskPlan
}

func NewMaskPlanCache() *MaskPlanCache {
	return &MaskPlanCache{plans: make(map[string]*MaskPlan)}
}

func (c *MaskPlanCache) Save(key string, plan *MaskPlan) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.plans[key] = plan
}

func (c *MaskPlanCache) Load(key string) *MaskPlan {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.plans[key]
}
