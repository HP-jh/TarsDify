package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateJSON(t *testing.T) {
	v := NewValidator()
	handler := v.ValidateJSON(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != 415 { t.Fatalf("expected 415, got %d", w.Code) }
	req2 := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler(w2, req2)
	if w2.Code != 200 { t.Fatalf("expected 200, got %d", w2.Code) }
}
