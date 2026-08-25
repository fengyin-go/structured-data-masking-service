package store

import (
	"datamasking/internal/model"
)

func (s *MemoryStore) CreateMaskRecord(r *model.MaskRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maskRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) GetMaskRecord(id string) (*model.MaskRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.maskRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListMaskRecords() []*model.MaskRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MaskRecord, 0, len(s.maskRecords))
	for _, r := range s.maskRecords {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) ListMaskRecordsByTask(taskID string) []*model.MaskRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MaskRecord, 0)
	for _, r := range s.maskRecords {
		if r.TaskID == taskID {
			list = append(list, r)
		}
	}
	return list
}

func (s *MemoryStore) DeleteMaskRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.maskRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.maskRecords, id)
	return nil
}
