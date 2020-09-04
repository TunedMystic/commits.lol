package config

import (
	"os"
	"path/filepath"
)

// Config contains all settings for the application.
type Config struct {
	BaseDir      string
	DatabaseName string
}

// GetConfig creates a Config type with settings from the environment.
func GetConfig() *Config {
	c := Config{}
	c.BaseDir = filepath.Dir(os.Getenv("BASE_DIR"))
	c.DatabaseName = os.Getenv("DATABASE_NAME")
	return &c
}
