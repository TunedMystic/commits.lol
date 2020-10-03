package db

import "github.com/tunedmystic/commits.lol/app/models"

// MockDB is an fake DB type that implements the Database interface.
// Used for testing.
type MockDB struct {
	RecentCommitsMock func() ([]*models.GitCommit, error)
}

func (m *MockDB) RecentCommits() ([]*models.GitCommit, error) {
	return m.RecentCommitsMock()
}

// Ensure the MockDB type satisfies the Database interface.
var _ Database = &MockDB{}
