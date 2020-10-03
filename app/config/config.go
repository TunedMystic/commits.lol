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
	if c.BaseDir == "" {
		panic("config: BaseDir not set")
	}
	// This could work to resolve the project root directory...
	// p, _ := filepath.Abs(c.BaseDir)
	// fmt.Printf("p:%v\n", p)

	c.DatabaseName = os.Getenv("DATABASE_NAME")
	if c.DatabaseName == "" {
		panic("config: DatabaseName not set")
	}
	return &c
}
