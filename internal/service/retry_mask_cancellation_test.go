package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"datamasking/internal/service"
	"datamasking/internal/store"
)

func TestRetryMaskStopsAfterRequestCancellation(t *testing.T) {
	client := store.NewRetryMaskClient(8 * time.Millisecond)
	runner := service.NewRetryMaskRunner(client)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Millisecond)
	defer cancel()
	if err := runner.Run(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("cancelled masking retry returned %v", err)
	}
	atReturn := client.Calls()
	time.Sleep(50 * time.Millisecond)
	if later := client.Calls(); later != atReturn {
		t.Errorf("masking calls kept growing after cancellation: at_return=%d later=%d", atReturn, later)
	}
}
