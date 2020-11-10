package db

import "github.com/tunedmystic/commits.lol/app/models"

// MockDB is an fake DB type that implements the Database interface.
// Used for testing.
type MockDB struct {
	AllBadWordsMock       func() (models.BadWords, error)
	AllGroupTermsMock     func() (models.GroupTerms, error)
	RandomSearchTermsMock func() models.SearchTerms

	AllCommitsMock           func() (models.GitCommits, error)
	UpdateCommitMock         func(commit *models.GitCommit) error
	RecentCommitsByGroupMock func(group string) (models.GitCommits, error)
	GetOrCreateUserMock      func(user *models.GitUser) error
	GetOrCreateRepoMock      func(repo *models.GitRepo) error
	GetOrCreateCommitMock    func(commit *models.GitCommit) error
}

// AllBadWords ...
func (m *MockDB) AllBadWords() (models.BadWords, error) {
	return m.AllBadWordsMock()
}

// AllGroupTerms ...
func (m *MockDB) AllGroupTerms() (models.GroupTerms, error) {
	return m.AllGroupTermsMock()
}

// RandomSearchTerms ...
func (m *MockDB) RandomSearchTerms() models.SearchTerms {
	return m.RandomSearchTermsMock()
}

// AllCommits ...
func (m *MockDB) AllCommits() (models.GitCommits, error) {
	return m.AllCommitsMock()
}

// UpdateCommit ...
func (m *MockDB) UpdateCommit(commit *models.GitCommit) error {
	return m.UpdateCommitMock(commit)
}

// RecentCommitsByGroup ...
func (m *MockDB) RecentCommitsByGroup(group string) (models.GitCommits, error) {
	return m.RecentCommitsByGroupMock(group)
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
