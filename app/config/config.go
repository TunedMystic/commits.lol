package config

import (
	"fmt"
	"path/filepath"
	"runtime"

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

// App stores the configuration for the application.
var App Config

// BasePath is the root directory of the project.
var BasePath string

func setBasePath(s *string) {
	_, b, _, _ := runtime.Caller(0)
	*s = filepath.Join(filepath.Dir(b), "../..")
}

func init() {
	// Load config variables from the environment.
	err := envconfig.Process("", &App)
	if err != nil {
		panic(fmt.Sprintf("config: %+v\n", err.Error()))
	}

	// Resolve basepath.
	setBasePath(&BasePath)
}
