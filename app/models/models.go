package models

import "time"

// GitSource ...
type GitSource struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// GitUser ...
type GitUser struct {
	ID        int    `db:"id"`
	SourceID  int    `db:"source_id"`
	Username  string `db:"username"`
	URL       string `db:"url"`
	AvatarURL string `db:"avatar_url"`

	Source GitSource `db:"source"`
}

// GitRepo ...
type GitRepo struct {
	ID          int    `db:"id"`
	SourceID    int    `db:"source_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	URL         string `db:"url"`

	Source *GitSource `db:"source"`
}

// GitCommit ...
type GitCommit struct {
	ID       int       `db:"id"`
	SourceID int       `db:"source_id"`
	AuthorID int       `db:"author_id"`
	RepoID   int       `db:"repo_id"`
	Message  string    `db:"message"`
	SHA      string    `db:"sha"`
	URL      string    `db:"url"`
	Date     time.Time `db:"date"`

	Source *GitSource `db:"source"`
	Author *GitUser   `db:"author"`
	Repo   *GitRepo   `db:"repo"`
}
