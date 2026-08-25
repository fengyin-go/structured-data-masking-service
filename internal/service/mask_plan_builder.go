package service

import (
	"fmt"

	"datamasking/internal/store"
)

type MaskPlanBuilder struct {
	cache *store.MaskPlanCache
}

func NewMaskPlanBuilder(cache *store.MaskPlanCache) *MaskPlanBuilder {
	return &MaskPlanBuilder{cache: cache}
}

func (b *MaskPlanBuilder) Build(key string, fail bool) (plan *store.MaskPlan, err error) {
	plan = &store.MaskPlan{}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("compile mask plan: %v", recovered)
		}
		b.cache.Save(key, plan)
	}()
	plan.Fields = append(plan.Fields, "email")
	if fail {
		panic("invalid policy")
	}
	plan.Ready = true
	return plan, nil
}
