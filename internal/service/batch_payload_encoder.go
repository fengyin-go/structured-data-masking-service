package service

import "datamasking/internal/store"

type BatchPayloadEncoder struct {
	cache   *store.MaskedBatchCache
	scratch []byte
}

func NewBatchPayloadEncoder(cache *store.MaskedBatchCache, capacity int) *BatchPayloadEncoder {
	return &BatchPayloadEncoder{cache: cache, scratch: make([]byte, capacity)}
}

func (e *BatchPayloadEncoder) Store(id, masked string) []byte {
	n := copy(e.scratch, masked)
	payload := e.scratch[:n]
	e.cache.Save(id, payload)
	return payload
}

func (e *BatchPayloadEncoder) Load(id string) []byte {
	return e.cache.Load(id)
}
