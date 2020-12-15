package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/tunedmystic/commits.lol/app/config"
	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/models"
	u "github.com/tunedmystic/commits.lol/app/utils"
)

func init() {
	// Switch to the base directory, so the `Server` type
	// can load templates correctly.
	os.Chdir(config.BasePath)
}

func mockGitCommits() models.GitCommits {
	return models.GitCommits{
		{Message: "Fixed a bug", MessageCensored: "Fixed a bug"},
		{Message: "Fixed another bug", MessageCensored: "Fixed another bug"},
		{Message: "Yolo changes", MessageCensored: "Yolo changes"},
	}
}

func Test_NewServer(t *testing.T) {
	NewServer(&db.MockDB{})
}

func TestRoutes(t *testing.T) {
	s := NewServer(&db.MockDB{})
	s.Routes()
}

func Test_IndexHandler_renders_index_page(t *testing.T) {
	mockDB := db.MockDB{
		RecentCommitsByGroupMock: func(group string) (models.GitCommits, error) {
			return mockGitCommits(), nil
		},
	}

	s := NewServer(&mockDB)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.IndexHandler).ServeHTTP(w, r)

	u.AssertEqual(t, w.Code, http.StatusOK)
}

func Test_IndexHandler_renders_commits_fragment(t *testing.T) {
	mockDB := db.MockDB{
		RecentCommitsByGroupMock: func(group string) (models.GitCommits, error) {
			return mockGitCommits(), nil
		},
	}

	s := NewServer(&mockDB)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Query params
	q := url.Values{}
	q.Add("fragment", "true")
	r.URL.RawQuery = q.Encode()

	// The handler will render a commits HTML fragment when
	// the `fragment` query param is true.
	http.HandlerFunc(s.IndexHandler).ServeHTTP(w, r)

	u.AssertEqual(t, w.Code, http.StatusOK)
}

func Test_IndexHandler_DB_error(t *testing.T) {
	mockDB := db.MockDB{
		RecentCommitsByGroupMock: func(group string) (models.GitCommits, error) {
			return models.GitCommits{}, errors.New("boom")
		},
	}

	s := NewServer(&mockDB)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(s.IndexHandler).ServeHTTP(w, r)

	u.AssertEqual(t, w.Code, http.StatusInternalServerError)

	// Check the body.
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	u.AssertEqual(t, string(body), "oopsie, something went horribly wrong\n")
}
