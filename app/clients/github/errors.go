package github

import (
	"encoding/json"
	"fmt"
)

// Errors raised by Commit validation.
const (
	ErrNoAuthor      ValidationError = "validate CommitItem: no author"
	ErrMessageLength ValidationError = "validate CommitItem: commit message too long"
	ErrMessageFormat ValidationError = "validate CommitItem: commit message has formatting issues"
)

// ValidationError is returned when a Commit is not valid.
type ValidationError string

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

func (err ValidationError) Error() string {
	return string(err)
}

func (e *APIError) Error() string {
	return fmt.Sprintf("github error %v: %v | URL: %v", e.StatusCode, e.Message, e.URL)
}
