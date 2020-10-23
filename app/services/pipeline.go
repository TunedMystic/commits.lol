package services

import (
	"fmt"

	"github.com/tunedmystic/commits.lol/app/clients/github"
	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/models"
)

// CommitPipeline is responsible for fetching commits
// concurrently and saving them to the database.
type CommitPipeline struct {
	db      db.Database
	client  *github.Client
	options github.CommitSearchOptions
	terms   []string
	done    chan bool
	source  *models.GitSource
}

// Commits creates and returns a commitPipeline type.
func Commits(db db.Database) *CommitPipeline {
	c := CommitPipeline{
		db:     db,
		client: github.NewClient(),
		done:   make(chan bool),
	}
	return &c
}

// WithTerms ...
func (c *CommitPipeline) WithTerms(terms ...string) *CommitPipeline {
	c.terms = append(c.terms, terms...)
	return c
}

// WithOptions ...
func (c *CommitPipeline) WithOptions(options github.CommitSearchOptions) *CommitPipeline {
	options.QueryText = ""
	options.Page = 1
	c.options = options
	return c
}

// Run ...
func (c *CommitPipeline) Run() error {
	fmt.Println("pipeline.Run")
	source, err := c.db.GetSource("github")

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("ok ...")
	c.source = source

	// Create goroutines for each term that we want to search for.
	for i, term := range c.terms {
		go c.fetch(i, term)
		fmt.Println("goroutine started")
	}

	// Wait for all goroutines to finish.
	for i := 0; i < len(c.terms); i++ {
		<-c.done
		fmt.Println("goroutine finished")
	}

	close(c.done)

	return nil
}

func (c *CommitPipeline) fetch(i int, term string) {
	fmt.Printf("[%v] searching for [%v]\n", i, term)
	options := c.options
	options.QueryText = term
	commitItems, err := c.client.CommitSearchPaginated(options)

	if err != nil {
		fmt.Printf("pipeline.fetch: %v\n", err)
		c.done <- true
		return
	}

	c.save(commitItems)
	c.done <- true
}

func (c *CommitPipeline) save(commitItems []github.CommitItem) {
	for _, commitItem := range commitItems {

		// Skip if commitItem is not valid.
		if err := commitItem.Validate(); err != nil {
			continue
		}

		// Transfer data from response struct into model struct.
		author := c.toAuthor(&commitItem)
		repo := c.toRepo(&commitItem)
		commit := c.toCommit(&commitItem)

		author.SourceID = c.source.ID
		repo.SourceID = c.source.ID
		commit.SourceID = c.source.ID

		// GetOrCreate Author
		if err := c.db.GetOrCreateUser(author); err != nil {
			fmt.Printf("pipeline.save:GetOrCreateUser: %v\n", err)
			continue
		}

		// GetOrCreate Repo
		if err := c.db.GetOrCreateRepo(repo); err != nil {
			fmt.Printf("pipeline.save:GetOrCreateRepo: %v\n", err)
			continue
		}

		// Create Commit with Author and Repo
		commit.AuthorID = author.ID
		commit.RepoID = repo.ID
		if err := c.db.GetOrCreateCommit(commit); err != nil {
			fmt.Printf("pipeline.save:GetOrCreateCommit: %v\n", err)
			continue
		}

		// fmt.Printf(">> %v | %v | %v\n\n", commit.Message, commit.ID, commit.URL)

		// fmt.Printf("%+v\n", author)
		// fmt.Printf("%+v\n", repo)
		// fmt.Printf("%+v\n", commit)
	}
}

func (c *CommitPipeline) toAuthor(item *github.CommitItem) *models.GitUser {
	return &models.GitUser{
		Username:  item.Author.Login,
		URL:       item.Author.URL,
		AvatarURL: item.Author.AvatarURL,
	}
}

func (c *CommitPipeline) toRepo(item *github.CommitItem) *models.GitRepo {
	return &models.GitRepo{
		Name:        item.Repo.Name,
		Description: item.Repo.Description,
		URL:         item.Repo.URL,
	}
}

func (c *CommitPipeline) toCommit(item *github.CommitItem) *models.GitCommit {
	return &models.GitCommit{
		Message: item.Commit.Message,
		SHA:     item.SHA,
		URL:     item.URL,
		Date:    item.Commit.Author.Date,
	}
}

/*
Usage:

db := db.NewSqliteDB()
options := github.CommitSearchOptions{}
err := Commits(db).WithOptions(options).WithTerms("hello").Run()

*/
