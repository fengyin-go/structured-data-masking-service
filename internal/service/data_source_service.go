package service

import (
	"sort"
	"time"

	"datamasking/internal/model"
	"datamasking/pkg/idgen"
)

func (s *Service) CreateDataSource(input model.DataSource) (*model.DataSource, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	d := &model.DataSource{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateDataSource(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDataSource(id string) (*model.DataSource, error) {
	return s.store.GetDataSource(id)
}

func (s *Service) ListDataSources(filter model.DataSourceFilter, page, size int) ([]*model.DataSource, int, error) {
	all := s.store.ListDataSources()
	matched := make([]*model.DataSource, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DataSource{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDataSource(id string, input model.DataSource) (*model.DataSource, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetDataSource(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Type = input.Type
	existing.Description = input.Description
	existing.UpdatedAt = time.Now().Unix()
	if err := s.store.UpdateDataSource(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteDataSource(id string) error {
	return s.store.DeleteDataSource(id)
}
