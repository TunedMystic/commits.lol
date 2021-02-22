package github

import (
	"regexp"
	"time"

	"github.com/tunedmystic/commits.lol/app/config"
)

// CommitMessagePattern is used to validate a commit message.
var CommitMessagePattern *regexp.Regexp

func init() {
	CommitMessagePattern = regexp.MustCompile(`^[a-zA-Z][\w '#%\.\!\:\-\)\(]+$`)
}

// RateLimitResponse ...
type RateLimitResponse struct {
	Resources struct {
		Core struct {
			Limit     int `json:"limit"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
		} `json:"core"`
		Search struct {
			Limit     int `json:"limit"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
		} `json:"search"`
	} `json:"resources"`
}

// CommitSearchResponse ...
type CommitSearchResponse struct {
	TotalCount  int          `json:"total_count"`
	CommitItems []CommitItem `json:"items"`
}

// IsEmpty checks if the response is empty.
func (resp CommitSearchResponse) IsEmpty() bool {
	return len(resp.CommitItems) == 0
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
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	Owner       User   `json:"owner"`
}

// Validate ...
func (c CommitItem) Validate() error {

	if c.Author == (User{}) {
		return ErrNoAuthor
	}

	if len(c.Commit.Message) > config.App.GithubCommitLength {
		return ErrMessageLength
	}

	if match := CommitMessagePattern.MatchString(c.Commit.Message); !match {
		return ErrMessageFormat
	}

	return nil
}
