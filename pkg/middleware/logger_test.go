package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	handler := Logger(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) }
	if w.Header().Get("X-Trace-ID") == "" { t.Fatal("expected trace-id") }
}
