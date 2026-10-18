package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBucket_Allow(t *testing.T) {
	b := NewBucket(3, time.Second)
	if !b.Allow() { t.Fatal("expected allow") }
	if !b.Allow() { t.Fatal("expected allow") }
	if !b.Allow() { t.Fatal("expected allow") }
	if b.Allow() { t.Fatal("expected deny") }
}

func TestHandler(t *testing.T) {
	b := NewBucket(2, time.Second)
	handler := Handler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, b)
	w1 := httptest.NewRecorder()
	handler(w1, httptest.NewRequest(http.MethodGet, "/", nil))
	if w1.Code != 200 { t.Fatalf("expected 200, got %d", w1.Code) }
	w2 := httptest.NewRecorder()
	handler(w2, httptest.NewRequest(http.MethodGet, "/", nil))
	if w2.Code != 200 { t.Fatalf("expected 200, got %d", w2.Code) }
	w3 := httptest.NewRecorder()
	handler(w3, httptest.NewRequest(http.MethodGet, "/", nil))
	if w3.Code != 429 { t.Fatalf("expected 429, got %d", w3.Code) }
}
