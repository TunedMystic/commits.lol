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

// RandomSearchTermsByRank returns a list of randomly selected terms of a specified rank.
func (s *SqliteDB) RandomSearchTermsByRank(amount, rank int) (models.SearchTerms, error) {
	terms := []models.SearchTerm{}

	query := `SELECT * FROM config_searchterm WHERE rank = ? ORDER BY random() LIMIT ?;`

	if err := s.DB.Select(&terms, query, rank, amount); err != nil {
		return nil, err
	}
	return models.SearchTerms(terms), nil
}

// RandomSearchTerms returns a list of randomly selected terms of predetermined rank.
func (s *SqliteDB) RandomSearchTerms() models.SearchTerms {
	values := []models.SearchTerm{}

	// Get terms of Rank 1.
	t1, err := s.RandomSearchTermsByRank(10, 1)
	if err != nil {
		fmt.Println(err)
	}
	values = append(values, t1...)

	// Get terms of Rank 2.
	t2, err := s.RandomSearchTermsByRank(4, 2)
	if err != nil {
		fmt.Println(err)
	}
	values = append(values, t2...)

	return models.SearchTerms(values)
}

// AllBadWords returns all the bad words.
func (s *SqliteDB) AllBadWords() (models.BadWords, error) {
	values := []models.BadWord{}

	query := `SELECT * FROM config_badword;`

	if err := s.DB.Select(&values, query); err != nil {
		return nil, err
	}

	return models.BadWords(values), nil
}

// AllGroupTerms returns all the group terms.
func (s *SqliteDB) AllGroupTerms() (models.GroupTerms, error) {
	values := []models.GroupTerm{}

	query := `SELECT * FROM config_groupterm;`

	if err := s.DB.Select(&values, query); err != nil {
		return nil, err
	}

	return models.GroupTerms(values), nil
}

// AllCommits returns all the commits.
func (s *SqliteDB) AllCommits() (models.GitCommits, error) {
	values := []models.GitCommit{}

	query := `SELECT * FROM git_commit;`

	if err := s.DB.Select(&values, query); err != nil {
		return nil, err
	}

	return values, nil
}

// UpdateCommit ...
func (s *SqliteDB) UpdateCommit(commit *models.GitCommit) error {
	query := `
		UPDATE git_commit
		SET
			source = :source, author_id = :author_id, repo_id = :repo_id,
			message = :message, message_censored = :message_censored,
			sha = :sha, url = :url, date = :date, created_at = :created_at,
			valid = :valid, groupname = :groupname,
			color_bg = :color_bg, color_fg = :color_fg
		WHERE id = :id;`

	_, err := s.DB.NamedExec(query, commit)

	if err != nil {
		return fmt.Errorf("error inserting commit: %v", err)
	}

	return nil
}

// RecentCommitsByGroup returns the most recent commits.
func (s *SqliteDB) RecentCommitsByGroup(group string) (models.GitCommits, error) {
	length := 20
	commits := make([]models.GitCommit, 0, length)

	query := `
		SELECT
			c.*,

			u.id AS "author.id",
			u.source AS "author.source",
			u.username AS "author.username",
			u.url AS "author.url",
			u.avatar_url AS "author.avatar_url"

		FROM git_commit c
		INNER JOIN git_user u on u.id = c.author_id
		WHERE (
			c.valid = TRUE AND
			c.date > datetime('now', $1) AND
			(
				($2 != '' AND c.groupname = $2)
				OR
				($2 = '' AND c.groupname IS NOT NULL)
			)
		)
		ORDER BY random()
		LIMIT $3;`

	rows, err := s.DB.Queryx(query, "-150 days", group, length)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.GitCommit
		err := rows.StructScan(&c)
		if err != nil {
			return nil, err
		}
		commits = append(commits, c)
	}

	return models.GitCommits(commits), nil
}

// CreateUser inserts a new User row and returns the ID.
func (s *SqliteDB) CreateUser(user *models.GitUser) error {
	query := `
		INSERT INTO git_user ("source", "username", "url", "avatar_url")
		VALUES (:source, :username, :url, :avatar_url);`

	row, err := s.DB.NamedExec(query, user)

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
	query := `SELECT id FROM git_user WHERE url = ?;`
	err := s.DB.QueryRow(query, user.URL).Scan(&user.ID)

	if err == sql.ErrNoRows {
		return s.CreateUser(user)
	}

	return err
}

// ------------------------------------------------------------------
// Get or Create Repo.

// CreateRepo inserts a new Repo row and returns the ID.
func (s *SqliteDB) CreateRepo(repo *models.GitRepo) error {
	query := `
		INSERT INTO git_repo ("source", "name", "description", "url")
		VALUES (:source, :name, :description, :url);`

	row, err := s.DB.NamedExec(query, repo)

	if err != nil {
		return fmt.Errorf("error inserting repo: %v", err)
	}

	id, _ := row.LastInsertId()

	repo.ID = int(id)
	return nil
}

// GetOrCreateRepo is a convenience method to get the provided Repo,
// or create it if it doesn't exist.
func (s *SqliteDB) GetOrCreateRepo(repo *models.GitRepo) error {
	query := `SELECT id FROM git_repo WHERE url = ?;`
	err := s.DB.QueryRow(query, repo.URL).Scan(&repo.ID)

	if err == sql.ErrNoRows {
		return s.CreateRepo(repo)
	}

	return err
}

// ------------------------------------------------------------------
// Get or Create Commit.

// CreateCommit inserts a new Commit row and returns the ID.
func (s *SqliteDB) CreateCommit(commit *models.GitCommit) error {
	query := `
		INSERT INTO git_commit (
			"source", "author_id", "repo_id", "message", "message_censored",
			"sha", "url", "date", "created_at", "valid", "groupname",
			"color_bg", "color_fg"
		)
		VALUES (
			:source, :author_id, :repo_id, :message, :message_censored,
			:sha, :url, :date, :created_at, :valid, :groupname,
			:color_bg, :color_fg
		);`

	row, err := s.DB.NamedExec(query, commit)

	if err != nil {
		return fmt.Errorf("error inserting commit: %v", err)
	}

	id, _ := row.LastInsertId()

	commit.ID = int(id)
	return nil
}

// GetOrCreateCommit is a convenience method to get the provided Commit,
// or create it if it doesn't exist.
func (s *SqliteDB) GetOrCreateCommit(commit *models.GitCommit) error {
	query := `SELECT id FROM git_commit WHERE author_id = ? AND message = ?;`
	err := s.DB.QueryRow(query, commit.AuthorID, commit.Message).Scan(&commit.ID)

	if err == sql.ErrNoRows {
		return s.CreateCommit(commit)
	}

	return err
}

// ------------------------------------------------------------------

// Ensure the SqliteDB type satisfies the Database interface.
var _ Database = &SqliteDB{}
