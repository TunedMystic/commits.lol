package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tunedmystic/commits.lol/app/config"
)

// CommitSearchResponse ...
type CommitSearchResponse struct {
	TotalCount  int          `json:"total_count"`
	CommitItems []CommitItem `json:"items"`
}

// CommitItem ...
type CommitItem struct {
	URL    string     `json:"html_url"`
	SHA    string     `json:"sha"`
	Commit Commit     `json:"commit"`
	Author User       `json:"author"`
	Repo   Repository `json:"repository"`
	Score  float64    `json:"score"`
}

// Commit ...
type Commit struct {
	Message string     `json:"message"`
	Author  AuthorInfo `json:"author"`
}

// AuthorInfo ...
type AuthorInfo struct {
	Date time.Time `json:"date"`
}

// User ...
type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"html_url"`
}

// Repository ...
type Repository struct {
	Name  string `json:"name"`
	URL   string `json:"html_url"`
	Owner User   `json:"owner"`
}

// APIError ...
type APIError struct {
	URL        string `json:"-"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

// Validate ...
func (c CommitItem) Validate() error {
	if c.Author == (User{}) {
		return errors.New("CommitItem: no author")
	}

	if len(c.Commit.Message) > config.App.GithubCommitLength {
		return errors.New("CommitItem: commit message too long")
	}

	if strings.ContainsAny(c.Commit.Message, "\n") {
		return errors.New("CommitItem: newline in commit message")
	}

	return nil
}

// NewAPIError ...
func NewAPIError(url string, data []byte, statusCode int) *APIError {
	e := APIError{URL: url, StatusCode: statusCode}
	if err := json.Unmarshal(data, &e); err != nil {
		e.Message = "not able to unmarshal error response"
	}
	return &e
}

// Error satisfies the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("github error %v: %v | URL: %v", e.StatusCode, e.Message, e.URL)
}
