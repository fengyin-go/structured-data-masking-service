package service

import (
	"fmt"
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/internal/store"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskTask(input model.MaskTask) (*model.MaskTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDataSource(input.DataSourceID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("data_source_id", "数据源不存在")
		}
		return nil, err
	}
	for _, rid := range input.RuleIDs {
		if _, err := s.store.GetMaskRule(rid); err != nil {
			if err == store.ErrNotFound {
				return nil, model.NewValidationError("rule_ids", fmt.Sprintf("规则 %s 不存在", rid))
			}
			return nil, err
		}
	}
	now := time.Now().Unix()
	t := &model.MaskTask{
		ID:           idgen.Hex(),
		DataSourceID: input.DataSourceID,
		RuleIDs:      input.RuleIDs,
		Status:       model.TaskStatusPending,
		Total:        input.Total,
		MaskedCount:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetMaskTask(id string) (*model.MaskTask, error) {
	return s.store.GetMaskTask(id)
}

func (s *Service) ListMaskTasks(filter model.MaskTaskFilter, page, size int) ([]*model.MaskTask, int, error) {
	all := s.store.ListMaskTasks()
	matched := make([]*model.MaskTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MaskTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMaskTaskStatus(id string, newStatus string) (*model.MaskTask, error) {
	if !validTaskStatuses[newStatus] {
		return nil, model.NewValidationError("status", "任务状态不合法")
	}
	t, err := s.store.GetMaskTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(t.Status, newStatus) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不允许从 %s 流转到 %s", t.Status, newStatus))
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now().Unix()
	if err := s.store.UpdateMaskTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteMaskTask(id string) error {
	return s.store.DeleteMaskTask(id)
}

var validTaskStatuses = map[string]bool{
	model.TaskStatusPending: true,
	model.TaskStatusRunning: true,
	model.TaskStatusDone:    true,
	model.TaskStatusFailed:  true,
}
