package db

import (
	"database/sql"
	"fmt"

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

// GetUserID retrieves a User ID, using the URL as the unique constraint.
func (s *SqliteDB) GetUserID(user *models.GitUser) error {
	sql := `SELECT id FROM git_user WHERE url = ?;`

	return s.DB.QueryRow(sql, user.URL).Scan(&user.ID)
}

// CreateUser inserts a new User row and returns the ID.
func (s *SqliteDB) CreateUser(user *models.GitUser) error {
	sql := `
		INSERT INTO git_user ("source_id", "username", "url", "avatar_url")
		VALUES (:source_id, :username, :url, :avatar_url);`

	row, err := s.DB.NamedExec(sql, user)

	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	id, _ := row.LastInsertId()

	user.ID = int(id)
	return nil
}

// GetOrCreateUser is a convenience method to get the provided User,
// or create it if it doesn't exist.
func (s *SqliteDB) GetOrCreateUser(user *models.GitUser) error {
	err := s.GetUserID(user)

	if err == sql.ErrNoRows {
		return s.CreateUser(user)
	}

	return err
}

// Ensure the SqliteDB type satisfies the Database interface.
var _ Database = &SqliteDB{}
