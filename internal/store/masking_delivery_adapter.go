package store

import (
	"errors"
	"sync"
)

type MaskingDeliveryAdapter struct {
	mu       sync.Mutex
	attempts int
	effects  int
}

func (a *MaskingDeliveryAdapter) Deliver(mode string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.attempts++
	a.effects++
	switch mode {
	case "reject":
		return errors.New("delivery rejected")
	case "temporary":
		if a.attempts == 1 {
			return errors.New("delivery unavailable")
		}
	}
	return nil
}

func (a *MaskingDeliveryAdapter) Counts() (int, int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.attempts, a.effects
}
