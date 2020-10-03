package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_GithubClient_OK(t *testing.T) {
	s := testServer(http.StatusOK, []byte(responseCommitSearch))
	defer s.Close()

	g := NewClient("some-api-key")
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, response.TotalCount, 1)
	u.AssertEqual(t, err, nil)
}

func Test_GithubClient_APIError(t *testing.T) {
	s := testServer(http.StatusUnprocessableEntity, []byte(responseValidationFailed))
	defer s.Close()

	g := NewClient("some-api-key")
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := "github error 422: Validation Failed | URL: http://127.0.0.1:59059/search/commits?q='fixed+a+bug'"

	response, err := g.CommitSearch(options)

	// How to detect whether a struct pointer is nil in golang?
	// Ref: https://stackoverflow.com/a/55900511
	u.AssertEqual(t, reflect.ValueOf(response).IsNil(), true)
	u.AssertEqual(t, err.Error(), expected)
}

func Test_GithubClient_empty_search_options(t *testing.T) {
	s := testServer(http.StatusOK, []byte(responseCommitSearch))
	defer s.Close()

	g := NewClient("some-api-key")
	g.baseURL = s.URL

	response, err := g.CommitSearch(CommitSearchOptions{})

	u.AssertEqual(t, reflect.ValueOf(response).IsNil(), true)
	u.AssertEqual(t, err.Error(), "no search options provided")
}

func Test_GithubClient_invalid_url(t *testing.T) {
	g := NewClient("some-api-key")
	g.baseURL = "1"

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := `error making request: Get "1/search/commits?q='fixed+a+bug'": unsupported protocol scheme ""`

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, reflect.ValueOf(response).IsNil(), true)
	u.AssertEqual(t, err.Error(), expected)
}

func Test_GithubClient_unmarshal_fail(t *testing.T) {
	s := testServer(http.StatusOK, []byte(`{"bad json"}`))
	defer s.Close()

	g := NewClient("some-api-key")
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := `not able to unmarshal response: invalid character '}' after object key`

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, reflect.ValueOf(response).IsNil(), true)
	u.AssertEqual(t, err.Error(), expected)
}

// ------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------

// testServer creates a testing server so we can mock our API responses.
func testServer(status int, data []byte) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(data))
	})
	srv := httptest.NewServer(handler)
	return srv
}
