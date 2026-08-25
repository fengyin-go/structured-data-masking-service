package service_test

import (
	"testing"

	"datamasking/internal/service"
	"datamasking/internal/store"
)

func didPanic(fn func()) (panicked bool) {
	defer func() { panicked = recover() != nil }()
	fn()
	return false
}

func TestPolicyZeroValuesRemainUsable(t *testing.T) {
	var registry store.PolicyRegistry
	if didPanic(func() { registry.Put("email", "partial") }) {
		t.Errorf("zero-value policy registry panicked on first write")
	}
	if got := registry.Get("email"); got != "partial" {
		t.Errorf("zero-value policy registry lost first policy: %q", got)
	}

	resolver := service.NewPolicyResolver(nil)
	if didPanic(func() { _ = resolver.Resolve("partial") }) {
		t.Errorf("missing optional policy validator caused a panic")
	}
}
