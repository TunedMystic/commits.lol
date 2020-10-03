package github

import (
	"encoding/json"
	"fmt"
	"time"
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
	Author Author     `json:"author"`
	Repo   Repository `json:"repository"`
	Score  float64    `json:"score"`
}

// Commit ...
type Commit struct {
	Message string     `json:"message"`
	Author  AuthorDate `json:"author"`
}

// AuthorDate ...
type AuthorDate struct {
	Date time.Time `json:"date"`
}

// Author ...
type Author struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"html_url"`
}

// Repository ...
type Repository struct {
	Name  string          `json:"name"`
	URL   string          `json:"html_url"`
	Owner RepositoryOwner `json:"owner"`
}

// RepositoryOwner ...
type RepositoryOwner struct {
	Login string `json:"login"`
}

// APIError ...
type APIError struct {
	URL        string `json:"-"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
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
