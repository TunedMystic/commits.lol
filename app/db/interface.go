package db

import "github.com/tunedmystic/commits.lol/app/models"

// Database defines the behavior for the application's database.
type Database interface {
	AllBadWords() (models.BadWords, error)
	AllGroupTerms() (models.GroupTerms, error)
	RandomSearchTerms() (models.SearchTerms, error)

	AllCommits() (models.GitCommits, error)
	UpdateCommit(commit *models.GitCommit) error
	RecentCommitsByGroup(group string) (models.GitCommits, error)
	GetOrCreateUser(user *models.GitUser) error
	GetOrCreateRepo(repo *models.GitRepo) error
	GetOrCreateCommit(commit *models.GitCommit) error
}
