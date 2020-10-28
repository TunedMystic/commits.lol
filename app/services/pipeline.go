package services

import (
	"fmt"
	"time"

	"github.com/tunedmystic/commits.lol/app/clients/github"
	"github.com/tunedmystic/commits.lol/app/config"
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
	jobs    chan string
	done    chan bool
	now     time.Time
}

// Commits creates and returns a CommitPipeline type.
func Commits(db db.Database) *CommitPipeline {
	c := CommitPipeline{
		db:     db,
		client: github.NewClient(),
		jobs:   make(chan string),
		done:   make(chan bool),
		now:    time.Now().UTC(),
	}

	return &c
}

// WithRandomTerms adds random terms to the pipeline.
func (c *CommitPipeline) WithRandomTerms() *CommitPipeline {
	c.terms = []string{}
	randomTerms := c.db.RandomTerms()
	c.terms = append(c.terms, randomTerms.Strings()...)
	return c
}

// WithTerms adds the provided terms to the pipeline.
func (c *CommitPipeline) WithTerms(terms ...string) *CommitPipeline {
	c.terms = []string{}
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
func (c *CommitPipeline) Run() {
	fmt.Println("pipeline.Run")

	// Exit if there are no terms.
	if len(c.terms) == 0 {
		fmt.Println("no terms in pipeline. exiting.")
		return
	}

	// Start the workers.
	for i := 0; i < config.WorkerSize; i++ {
		go c.worker(i)
	}

	// Write jobs to the jobs channel.
	go c.writeJobs()

	// Wait for all goroutines to finish.
	for i := 0; i < len(c.terms); i++ {
		<-c.done
		fmt.Println("goroutine finished")
	}

	close(c.done)
}

// writeJobs sends jobs to the jobs channel and then closes the channel.
func (c *CommitPipeline) writeJobs() {
	for _, term := range c.terms {
		c.jobs <- term
	}
	close(c.jobs)
}

// worker consumes jobs from the jobs channel, and executes the work.
func (c *CommitPipeline) worker(ID int) {
	fmt.Printf("worker [%v] started \n", ID)
	for term := range c.jobs {
		// Copy the options, and add the term.
		options := c.options
		options.QueryText = term

		// Perform the commit search.
		commitItems, err := c.client.CommitSearchPaginated(options)

		if err != nil {
			fmt.Printf("pipeline.fetch: %v\n", err)
			c.done <- true
			continue
		}

		// Save commitItems to the database.
		for _, commitItem := range commitItems {
			c.save(commitItem)
		}

		c.done <- true
	}
	fmt.Printf("worker [%v] done\n", ID)
}

func (c *CommitPipeline) save(commitItem github.CommitItem) error {
	// Skip if commitItem is not valid.
	if err := commitItem.Validate(); err != nil {
		return err
	}

	// Transfer data from response struct into model struct.
	author := c.toAuthor(commitItem)
	repo := c.toRepo(commitItem)
	commit := c.toCommit(commitItem)

	// GetOrCreate Author
	if err := c.db.GetOrCreateUser(&author); err != nil {
		return fmt.Errorf("pipeline.save:GetOrCreateUser: %v", err)
	}

	// GetOrCreate Repo
	if err := c.db.GetOrCreateRepo(&repo); err != nil {
		return fmt.Errorf("pipeline.save:GetOrCreateRepo: %v", err)
	}

	// Create Commit with Author and Repo
	commit.AuthorID = author.ID
	commit.RepoID = repo.ID
	if err := c.db.GetOrCreateCommit(&commit); err != nil {
		return fmt.Errorf("pipeline.save:GetOrCreateCommit: %v", err)
	}

	// fmt.Printf(">> %v | %v | %v\n\n", commit.Message, commit.ID, commit.URL)
	return nil
}

func (c *CommitPipeline) toAuthor(item github.CommitItem) models.GitUser {
	return models.GitUser{
		Source:    config.SourceGithub,
		Username:  item.Author.Login,
		URL:       item.Author.URL,
		AvatarURL: item.Author.AvatarURL,
	}
}

func (c *CommitPipeline) toRepo(item github.CommitItem) models.GitRepo {
	return models.GitRepo{
		Source:      config.SourceGithub,
		Name:        item.Repo.Name,
		Description: item.Repo.Description,
		URL:         item.Repo.URL,
	}
}

func (c *CommitPipeline) toCommit(item github.CommitItem) models.GitCommit {
	return models.GitCommit{
		Source:    config.SourceGithub,
		Message:   item.Commit.Message,
		SHA:       item.SHA,
		URL:       item.URL,
		Date:      item.Commit.Author.Date,
		CreatedAt: c.now,
		Valid:     true,
	}
}
