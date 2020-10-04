package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite

	"github.com/tunedmystic/commits.lol/app/models"
)

// SqliteDB is an sqlite-backed type that implements the Database interface.
type SqliteDB struct {
	DB *sqlx.DB
}

// NewSqliteDB connects to the database, and returns a new *SqliteDB type.
func NewSqliteDB(name string) *SqliteDB {
	sdb := SqliteDB{
		DB: sqlx.MustConnect("sqlite3", name),
	}
	return &sdb
}

// RecentCommits returns the most recent commits.
func (s *SqliteDB) RecentCommits() ([]*models.GitCommit, error) {
	commits := []*models.GitCommit{}

	sql := `
		SELECT
			c.*,

			u.id AS "author.id",
			u.source_id AS "author.source_id",
			u.username AS "author.username",
			u.url AS "author.url",
			u.avatar_url AS "author.avatar_url"

		FROM git_commit c
		INNER JOIN git_user u on u.id = c.author_id
		WHERE c.date > '2015-09-02';`

	if err := s.DB.Select(&commits, sql); err != nil {
		return nil, err
	}

	return commits, nil
}

// Ensure the SqliteDB type satisfies the Database interface.
var _ Database = &SqliteDB{}
