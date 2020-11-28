package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config contains all settings for the application.
type Config struct {
	Environment        string `split_words:"true" default:"dev"`
	BaseURL            string `split_words:"true" required:"true"`
	Port               int    `split_words:"true" required:"true"`
	DatabaseName       string `split_words:"true" required:"true"`
	GithubAPIKey       string `split_words:"true" required:"true"`
	GithubMaxFetch     int    `split_words:"true" default:"50"`
	GithubCommitLength int    `split_words:"true" default:"45"`
	LogLevel           string `split_words:"true" default:"INFO"`
	SentryDSN          string `split_words:"true"`
	GoatcounterUser    string `split_words:"true"`
}

// SourceGithub is an enum for the Github source.
const SourceGithub int = 1

// WorkerSize defines the amount of goroutines to spawn when running background tasks.
const WorkerSize int = 4

// GoatCounterScript initializes analytics for the site.
const GoatCounterScript string = `<script data-goatcounter="https://%s.goatcounter.com/count" async src="//gc.zgo.at/count.js"></script>`

// App stores the configuration for the application.
var App Config

func init() {
	// Load config variables from the environment.
	err := envconfig.Process("", &App)
	if err != nil {
		panic(fmt.Sprintf("config: %+v\n", err.Error()))
	}
}
