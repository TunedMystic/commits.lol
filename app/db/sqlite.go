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
func NewSqliteDB(name string) SqliteDB {
	sdb := SqliteDB{
		DB: sqlx.MustConnect("sqlite3", name),
	}
	return sdb
}

// Close the database connection.
func (s *SqliteDB) Close() {
	s.DB.Close()
}

// ------------------------------------------------------------------
// Methods to modify config-related tables (BadWWord, GroupTerm, SearchTerm)

// AllBadWords returns all the bad words.
func (s *SqliteDB) AllBadWords() (models.BadWords, error) {
	values := []models.BadWord{}

	if err := s.DB.Select(&values, `SELECT * FROM config_badword;`); err != nil {
		return nil, err
	}

	return models.BadWords(values), nil
}

// AllGroupTerms returns all the group terms.
func (s *SqliteDB) AllGroupTerms() (models.GroupTerms, error) {
	values := []models.GroupTerm{}

	if err := s.DB.Select(&values, `SELECT * FROM config_groupterm;`); err != nil {
		return nil, err
	}

	return models.GroupTerms(values), nil
}

// randomSearchTermsByRank returns a list of randomly selected terms of a specified rank.
func (s *SqliteDB) randomSearchTermsByRank(rank, amount int) (models.SearchTerms, error) {
	terms := make(models.SearchTerms, 0, amount)
	query := `SELECT * FROM config_searchterm WHERE rank = ? ORDER BY random() LIMIT ?;`

	if err := s.DB.Select(&terms, query, rank, amount); err != nil {
		return nil, err
	}
	return terms, nil
}

// RandomSearchTerms returns a list of randomly selected terms of predetermined rank.
func (s *SqliteDB) RandomSearchTerms() (models.SearchTerms, error) {
	rank1Amount := 8
	rank2Amount := 4
	rank3Amount := 4
	rank4Amount := 2
	totalTerms := rank1Amount + rank2Amount + rank3Amount + rank4Amount
	terms := make(models.SearchTerms, 0, totalTerms)

	// Get terms of Rank 1.
	t1, err := s.randomSearchTermsByRank(1, rank1Amount)
	if err != nil {
		return nil, fmt.Errorf("db:RandomSearchTerms: %v", err)
	}
	terms = append(terms, t1...)

	// Get terms of Rank 2.
	t2, err := s.randomSearchTermsByRank(2, rank2Amount)
	if err != nil {
		return nil, fmt.Errorf("db:RandomSearchTerms: %v", err)
	}
	terms = append(terms, t2...)

	// Get terms of Rank 3.
	t3, err := s.randomSearchTermsByRank(3, rank3Amount)
	if err != nil {
		return nil, fmt.Errorf("db:RandomSearchTerms: %v", err)
	}
	terms = append(terms, t3...)

	// Get terms of Rank 4.
	t4, err := s.randomSearchTermsByRank(4, rank4Amount)
	if err != nil {
		return nil, fmt.Errorf("db:RandomSearchTerms: %v", err)
	}
	terms = append(terms, t4...)

	return terms, nil
}

// ------------------------------------------------------------------
// Methods to modify git-related tables (GitCommit, GitRepo, GitUser)

// AllCommits returns all the commits.
func (s *SqliteDB) AllCommits() (models.GitCommits, error) {
	commits := make(models.GitCommits, 0, 1000)

	if err := s.DB.Select(&commits, `SELECT * FROM git_commit;`); err != nil {
		return nil, err
	}

	return commits, nil
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
	length := 33
	commits := make(models.GitCommits, 0, length)

	query := `
		SELECT
			c.id,
			c.author_id,
			c.repo_id,
			c.message,
			c.message_censored,
			c.url,
			c.color_bg,
			c.color_fg,

			u.id AS "author.id",
			u.url AS "author.url",
			u.avatar_url AS "author.avatar_url"

		FROM git_commit c
		INNER JOIN git_user u on u.id = c.author_id
		WHERE (
			c.date > datetime('now', '-14 days') AND
			c.valid = TRUE AND
			c.id % abs(random() % 10) = 0 AND
			(
				($1 != '' AND c.groupname = $1)
				OR
				($1 = '' AND c.groupname IS NOT NULL)
			)
		)
		ORDER BY random()
		LIMIT $2;`

	rows, err := s.DB.Queryx(query, group, length)

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

	return commits, nil
}

// createUser inserts a new User row and returns the ID.
func (s *SqliteDB) createUser(user *models.GitUser) error {
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
		return s.createUser(user)
	}

	return err
}

// createRepo inserts a new Repo row and returns the ID.
func (s *SqliteDB) createRepo(repo *models.GitRepo) error {
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
		return s.createRepo(repo)
	}

	return err
}

// createCommit inserts a new Commit row and returns the ID.
func (s *SqliteDB) createCommit(commit *models.GitCommit) error {
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
		return s.createCommit(commit)
	}

	return err
}

// ------------------------------------------------------------------

// Ensure the SqliteDB type satisfies the Database interface.
var _ Database = &SqliteDB{}
