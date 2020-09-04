package db

import "github.com/tunedmystic/commits.lol/app/models"

// Database defines the behavior for the application's database.
type Database interface {
	RecentCommits() ([]*models.GitCommit, error)
}
