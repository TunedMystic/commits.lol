package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_CommitSearch_OK(t *testing.T) {
	s := testServer(http.StatusOK, []byte(responseCommitSearch))
	defer s.Close()

	g := NewClient()
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, response.TotalCount, 1)
	u.AssertEqual(t, err, nil)
}

func Test_CommitSearch_APIError(t *testing.T) {
	s := testServer(http.StatusUnprocessableEntity, []byte(responseValidationFailed))
	defer s.Close()

	g := NewClient()
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := fmt.Sprintf("github error 422: Validation Failed | URL: %v/search/commits?q='fixed+a+bug'", s.URL)

	response, err := g.CommitSearch(options)

	// How to detect whether a struct pointer is nil in golang?
	// Ref: https://stackoverflow.com/a/55900511
	// u.AssertEqual(t, reflect.ValueOf(response).IsNil(), true)

	u.AssertEqual(t, response.IsEmpty(), true)
	u.AssertEqual(t, err.Error(), expected)
}

func Test_CommitSearch_empty_search_options(t *testing.T) {
	s := testServer(http.StatusOK, []byte(responseCommitSearch))
	defer s.Close()

	g := NewClient()
	g.baseURL = s.URL

	response, err := g.CommitSearch(CommitSearchOptions{})

	u.AssertEqual(t, response.IsEmpty(), true)
	u.AssertEqual(t, err.Error(), "no search options provided")
}

func Test_CommitSearch_invalid_url(t *testing.T) {
	g := NewClient()
	g.baseURL = "1"

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := `error making request: Get "1/search/commits?q='fixed+a+bug'": unsupported protocol scheme ""`

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, response.IsEmpty(), true)
	u.AssertEqual(t, err.Error(), expected)
}

func Test_CommitSearch_unmarshal_fail(t *testing.T) {
	s := testServer(http.StatusOK, []byte(`{"bad json"}`))
	defer s.Close()

	g := NewClient()
	g.baseURL = s.URL

	options := CommitSearchOptions{QueryText: "fixed a bug"}
	expected := `not able to unmarshal response: invalid character '}' after object key`

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, response.IsEmpty(), true)
	u.AssertEqual(t, err.Error(), expected)
}

func Test_CommitSearchPaginated(t *testing.T) {
	s := testServer(http.StatusOK, []byte(responseCommitSearchMany))
	defer s.Close()

	g := NewClient()
	g.baseURL = s.URL
	g.maxFetch = 10

	commitItems, err := g.CommitSearchPaginated(CommitSearchOptions{QueryText: "fixed a bug"})

	u.AssertEqual(t, len(commitItems), 10)
	u.AssertEqual(t, err, nil)
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
