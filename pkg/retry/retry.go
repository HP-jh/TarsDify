package retry

import (
	"context"
	"fmt"
	"time"
)

type Strategy struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func Default() *Strategy {
	return &Strategy{MaxRetries: 3, BaseDelay: 200 * time.Millisecond, MaxDelay: 2 * time.Second}
}

func (s *Strategy) Do(ctx context.Context, fn func() error) error {
	var err error
	for i := 0; i <= s.MaxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i == s.MaxRetries { break }
		delay := s.BaseDelay * (1 << i)
		if delay > s.MaxDelay { delay = s.MaxDelay }
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return fmt.Errorf("retry exhausted (%d): %w", s.MaxRetries, err)
}
