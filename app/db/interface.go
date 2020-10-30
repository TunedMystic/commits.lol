package db

import "github.com/tunedmystic/commits.lol/app/models"

// Database defines the behavior for the application's database.
type Database interface {
	RandomTermsByRank(amount, rank int) (models.Terms, error)
	RandomTerms() models.Terms
	AllCommits() ([]*models.GitCommit, error)
	UpdateCommit(commit *models.GitCommit) error
	RecentCommits() (models.GitCommits, error)
	GetOrCreateUser(user *models.GitUser) error
	GetOrCreateRepo(repo *models.GitRepo) error
	GetOrCreateCommit(commit *models.GitCommit) error
}
