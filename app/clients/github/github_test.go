package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_GithubClient(t *testing.T) {
	s := testServer()
	defer s.Close()

	g := NewClient("some-api-key")
	g.SetBaseURL(s.URL)
	// g.SetBaseURL(fmt.Sprint(1))
	// g.SetBaseURL(fmt.Sprint(0x7f))

	options := CommitSearchOptions{QueryText: "lol"}

	response, err := g.CommitSearch(options)

	u.AssertEqual(t, err, nil)
	u.AssertEqual(t, response.TotalCount, 1)
}

func testServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseCommitSearch))
	})
	srv := httptest.NewServer(handler)
	return srv
}
