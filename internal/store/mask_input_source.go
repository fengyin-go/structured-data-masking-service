package store

import (
	"context"
	"errors"
)

type MaskInputSource struct{}

func (MaskInputSource) Stream(ctx context.Context, inputs []string) (<-chan string, <-chan error) {
	values := make(chan string)
	errs := make(chan error)
	go func() {
		for _, input := range inputs {
			if input == "bad" {
				errs <- errors.New("invalid masking input")
				return
			}
			values <- input
		}
		close(values)
	}()
	return values, errs
}
