package service

import (
	"context"

	"datamasking/internal/store"
)

type MaskInputCollector struct {
	source store.MaskInputSource
}

func NewMaskInputCollector(source store.MaskInputSource) *MaskInputCollector {
	return &MaskInputCollector{source: source}
}

func (c *MaskInputCollector) Collect(ctx context.Context, inputs []string) ([]string, error) {
	values, _ := c.source.Stream(ctx, inputs)
	result := make([]string, 0, len(inputs))
	for value := range values {
		result = append(result, value)
	}
	return result, nil
}
