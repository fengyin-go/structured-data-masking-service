package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateMaskRule(input model.MaskRule) (*model.MaskRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	r := &model.MaskRule{
		ID:           idgen.Hex(),
		Name:         input.Name,
		FieldName:    input.FieldName,
		RuleType:     input.RuleType,
		KeepPrefix:   input.KeepPrefix,
		KeepSuffix:   input.KeepSuffix,
		MaskChar:     input.MaskChar,
		RegexPattern: input.RegexPattern,
		Enabled:      input.Enabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreateMaskRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetMaskRule(id string) (*model.MaskRule, error) {
	return s.store.GetMaskRule(id)
}

func (s *Service) ListMaskRules(filter model.MaskRuleFilter, page, size int) ([]*model.MaskRule, int, error) {
	all := s.store.ListMaskRules()
	matched := make([]*model.MaskRule, 0, len(all))
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
		return []*model.MaskRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMaskRule(id string, input model.MaskRule) (*model.MaskRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetMaskRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.FieldName = input.FieldName
	existing.RuleType = input.RuleType
	existing.KeepPrefix = input.KeepPrefix
	existing.KeepSuffix = input.KeepSuffix
	existing.MaskChar = input.MaskChar
	existing.RegexPattern = input.RegexPattern
	existing.Enabled = input.Enabled
	existing.UpdatedAt = time.Now().Unix()
	if err := s.store.UpdateMaskRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteMaskRule(id string) error {
	return s.store.DeleteMaskRule(id)
}

func (s *Service) BatchCreateMaskRules(inputs []model.MaskRule) ([]*model.MaskRule, error) {
	results := make([]*model.MaskRule, 0, len(inputs))
	for _, input := range inputs {
		r, err := s.CreateMaskRule(input)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
