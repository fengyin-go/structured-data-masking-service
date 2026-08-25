package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskRecord(input model.MaskRecord) (*model.MaskRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMaskTask(input.TaskID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("task_id", "任务不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetMaskRule(input.RuleID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("rule_id", "规则不存在")
		}
		return nil, err
	}
	r := &model.MaskRecord{
		ID:        idgen.Hex(),
		TaskID:    input.TaskID,
		RuleID:    input.RuleID,
		Original:  input.Original,
		Masked:    input.Masked,
		CreatedAt: time.Now().Unix(),
	}
	if err := s.store.CreateMaskRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetMaskRecord(id string) (*model.MaskRecord, error) {
	return s.store.GetMaskRecord(id)
}

func (s *Service) ListMaskRecords(filter model.MaskRecordFilter, page, size int) ([]*model.MaskRecord, int, error) {
	all := s.store.ListMaskRecords()
	matched := make([]*model.MaskRecord, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) ListMaskRecordsByTask(taskID string, page, size int) ([]*model.MaskRecord, int, error) {
	all := s.store.ListMaskRecordsByTask(taskID)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt > all[j].CreatedAt
	})
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *Service) DeleteMaskRecord(id string) error {
	return s.store.DeleteMaskRecord(id)
}
