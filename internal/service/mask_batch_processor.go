package service

import (
	"errors"

	"datamasking/internal/store"
)

type MaskBatchProcessor struct {
	pool *store.MaskResourcePool
}

func NewMaskBatchProcessor(pool *store.MaskResourcePool) *MaskBatchProcessor {
	return &MaskBatchProcessor{pool: pool}
}

func (p *MaskBatchProcessor) Process(items []string) (err error) {
	for _, item := range items {
		session, openErr := p.pool.Open()
		if openErr != nil {
			return openErr
		}
		defer session.Close()
		defer func() { err = session.Finish(err) }()
		if item == "invalid" {
			return errors.New("invalid masking field")
		}
	}
	return nil
}
