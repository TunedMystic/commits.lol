package config

import (
	"os"
	"strconv"
)

// Config contains all settings for the application.
type Config struct {
	DatabaseName       string
	GithubAPIKey       string
	GithubMaxFetch     int
	GithubCommitLength int
}

// GetConfig creates a Config type with settings from the environment.
func GetConfig() *Config {
	var err error
	c := Config{}

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

	return &c
}
