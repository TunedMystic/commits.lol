package config

import (
	"os"
	"strconv"
)

// SourceGithub is an enum for the Github source.
const SourceGithub int = 1

// WorkerSize defines the amount of goroutines to spawn when running background tasks.
const WorkerSize int = 4

// App stores the configuration for the application.
var App Config

func init() {
	App = GetConfig()
}

// Config contains all settings for the application.
type Config struct {
	BaseURL            string
	DatabaseName       string
	GithubAPIKey       string
	GithubMaxFetch     int
	GithubCommitLength int
}

// GetConfig creates a Config type with settings from the environment.
func GetConfig() Config {
	var err error
	c := Config{}

	c.BaseURL = os.Getenv("BASE_URL")
	if c.BaseURL == "" {
		panic("config: BaseURL not set")
	}

	c.DatabaseName = os.Getenv("DATABASE_NAME")
	if c.DatabaseName == "" {
		panic("config: DatabaseName not set")
	}

	c.GithubAPIKey = os.Getenv("GITHUB_TOKEN")
	if c.GithubAPIKey == "" {
		panic("config: GithubAPIKey not set")
	}

	maxFetch := os.Getenv("GITHUB_MAX_FETCH")
	if maxFetch == "" {
		maxFetch = "50"
	}
	c.GithubMaxFetch, err = strconv.Atoi(maxFetch)
	if err != nil {
		panic("config: GithubMaxFetch error when parsing")
	}

	commitLength := os.Getenv("GITHUB_COMMIT_LENGTH")
	if commitLength == "" {
		commitLength = "55"
	}
	c.GithubCommitLength, err = strconv.Atoi(commitLength)
	if err != nil {
		panic("config: GithubCommitLength error when parsing")
	}

	return c
}
