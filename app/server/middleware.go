package server

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// CacheControl ...
func CacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=300") // 5 minutes
		h.ServeHTTP(w, r)
	})
}

// StatusRecorder allows us to capture the response status code.
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

// WriteHeader ...
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logging ...
func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		h.ServeHTTP(rec, r)
		zap.L().Info("[" + strconv.Itoa(rec.Status) + "] " + r.Method + " " + r.URL.Path + " " + r.URL.RawQuery)
	})
}
