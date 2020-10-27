package db

import "github.com/tunedmystic/commits.lol/app/models"

// Database defines the behavior for the application's database.
type Database interface {
	RandomTerms(size, rank int) (models.Terms, error)
	RecentCommits() ([]*models.GitCommit, error)
	GetOrCreateUser(user *models.GitUser) error
	GetOrCreateRepo(repo *models.GitRepo) error
	GetOrCreateCommit(commit *models.GitCommit) error
}
