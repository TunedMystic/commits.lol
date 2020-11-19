package server

import (
	"net/http"

	"go.uber.org/zap"
)

// CacheControl ...
func CacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=300") // 5 minutes
		h.ServeHTTP(w, r)
	})
}

// Logging ...
func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		zap.L().Info(r.Method + " " + r.URL.Path + " " + r.URL.RawQuery)
	})
}
