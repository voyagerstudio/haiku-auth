package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all env var config required by haiku-auth
type Config struct {
	API *APIConfig
	DB  *DBConfig
}

// DefaultConfig returns sane defaults for commonly used deployment envs
func DefaultConfig() (*Config, error) {
	// This binds the environment to viper for reading env var config
	viper.AutomaticEnv()

	apiConfig, err := NewAPIConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading api config: %v", err)
	}

	dbConfig, err := NewDBConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading db config: %v", err)
	}

	c := &Config{
		API: apiConfig,
		DB:  dbConfig,
	}
	return c, nil
}
