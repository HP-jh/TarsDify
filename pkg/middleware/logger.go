package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(b)
	rr.size += n
	return n, err
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := r.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("%d", start.UnixNano())
		}
		w.Header().Set("X-Trace-ID", traceID)
		rr := &responseRecorder{ResponseWriter: w, status: 200}
		next(rr, r)
		log.Printf("[%s] %s %s %d %dB %v", traceID, r.Method, r.URL.Path, rr.status, rr.size, time.Since(start))
	}
}
