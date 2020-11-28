package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_CacheControl(t *testing.T) {
	// Dummy handler that will be wrapped with the middleware.
	handlerReached := false
	h := func(w http.ResponseWriter, r *http.Request) {
		handlerReached = true
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	CacheControl(http.HandlerFunc(h)).ServeHTTP(w, r)

	u.AssertEqual(t, w.Header().Get("Cache-Control"), "max-age=300")
	u.AssertEqual(t, handlerReached, true)
}

func Test_Logging(t *testing.T) {
	// Dummy handler that will be wrapped with the middleware.
	handlerReached := false
	h := func(w http.ResponseWriter, r *http.Request) {
		handlerReached = true
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Logging(http.HandlerFunc(h)).ServeHTTP(w, r)

	u.AssertEqual(t, handlerReached, true)
}
