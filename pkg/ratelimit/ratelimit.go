package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type Bucket struct {
	tokens   int
	capacity int
	last     time.Time
	mu       sync.Mutex
}

func NewBucket(capacity int, fillInterval time.Duration) *Bucket {
	return &Bucket{tokens: capacity, capacity: capacity, last: time.Now()}
}

func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last)
	refill := int(elapsed.Seconds())
	if refill > 0 {
		b.tokens += refill
		if b.tokens > b.capacity { b.tokens = b.capacity }
		b.last = now
	}
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func Handler(next http.HandlerFunc, bucket *Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !bucket.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}
		next(w, r)
	}
}
