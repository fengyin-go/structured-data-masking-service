package store

import (
	"datamasking/internal/model"
)

func (s *MemoryStore) CreateDataSource(d *model.DataSource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.dataSources {
		if exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.dataSources[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDataSource(id string) (*model.DataSource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dataSources[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDataSources() []*model.DataSource {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DataSource, 0, len(s.dataSources))
	for _, d := range s.dataSources {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDataSource(d *model.DataSource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataSources[d.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.dataSources {
		if exist.ID != d.ID && exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.dataSources[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDataSource(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dataSources[id]; !ok {
		return ErrNotFound
	}
	delete(s.dataSources, id)
	return nil
}
