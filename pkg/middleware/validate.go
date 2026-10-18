package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Validator struct {
	MaxBodySize int64
}

func NewValidator() *Validator {
	return &Validator{MaxBodySize: 1024 * 1024}
}

func (v *Validator) LimitBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, v.MaxBodySize)
		next(w, r)
	}
}

func (v *Validator) ValidateJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { next(w, r); return }
		ct := strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(ct, "application/json") && r.ContentLength > 0 {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "Content-Type must be application/json"})
			return
		}
		next(w, r)
	}
}
