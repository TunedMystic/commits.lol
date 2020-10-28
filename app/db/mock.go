package db

import "github.com/tunedmystic/commits.lol/app/models"

// MockDB is an fake DB type that implements the Database interface.
// Used for testing.
type MockDB struct {
	RandomTermsByRankMock func(amount, rank int) (models.Terms, error)
	RandomTermsMock       func() models.Terms
	RecentCommitsMock     func() (models.GitCommits, error)
	GetOrCreateUserMock   func(user *models.GitUser) error
	GetOrCreateRepoMock   func(repo *models.GitRepo) error
	GetOrCreateCommitMock func(commit *models.GitCommit) error
}

// RandomTermsByRank ...
func (m *MockDB) RandomTermsByRank(amount, rank int) (models.Terms, error) {
	return m.RandomTermsByRankMock(amount, rank)
}

// RandomTerms ...
func (m *MockDB) RandomTerms() models.Terms {
	return m.RandomTermsMock()
}

// RecentCommits ...
func (m *MockDB) RecentCommits() (models.GitCommits, error) {
	return m.RecentCommitsMock()
}

// GetOrCreateUser ...
func (m *MockDB) GetOrCreateUser(user *models.GitUser) error {
	return m.GetOrCreateUserMock(user)
}

// GetOrCreateRepo ...
func (m *MockDB) GetOrCreateRepo(repo *models.GitRepo) error {
	return m.GetOrCreateRepoMock(repo)
}

// GetOrCreateCommit ...
func (m *MockDB) GetOrCreateCommit(commit *models.GitCommit) error {
	return m.GetOrCreateCommitMock(commit)
}

// Ensure the MockDB type satisfies the Database interface.
var _ Database = &MockDB{}
