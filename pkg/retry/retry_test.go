package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_Success(t *testing.T) {
	s := &Strategy{MaxRetries: 2, BaseDelay: 10 * time.Millisecond, MaxDelay: 50 * time.Millisecond}
	calls := 0
	err := s.Do(context.Background(), func() error {
		calls++
		if calls < 2 { return errors.New("temp fail") }
		return nil
	})
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if calls != 2 { t.Fatalf("expected 2 calls, got %d", calls) }
}

func TestRetry_Exhausted(t *testing.T) {
	s := &Strategy{MaxRetries: 1, BaseDelay: 10 * time.Millisecond, MaxDelay: 50 * time.Millisecond}
	err := s.Do(context.Background(), func() error {
		return errors.New("always fail")
	})
	if err == nil { t.Fatal("expected error") }
}
