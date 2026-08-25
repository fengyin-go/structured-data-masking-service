package store

import (
	"errors"
	"sync"
)

type MaskResourcePool struct {
	mu        sync.Mutex
	limit     int
	open      int
	committed int
}

func NewMaskResourcePool(limit int) *MaskResourcePool {
	return &MaskResourcePool{limit: limit}
}

type MaskSession struct {
	pool   *MaskResourcePool
	closed bool
}

func (p *MaskResourcePool) Open() (*MaskSession, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.open >= p.limit {
		return nil, errors.New("mask resource pool exhausted")
	}
	p.open++
	return &MaskSession{pool: p}, nil
}

func (s *MaskSession) Finish(cause error) error {
	s.pool.mu.Lock()
	s.pool.committed++
	s.pool.mu.Unlock()
	return nil
}

func (s *MaskSession) Close() {
	if s.closed {
		return
	}
	s.closed = true
	s.pool.mu.Lock()
	s.pool.open--
	s.pool.mu.Unlock()
}

func (p *MaskResourcePool) Counts() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.open, p.committed
}
