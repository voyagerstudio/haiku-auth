package config

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DBConfig ...
type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Pass     string
}

// NewDBConfig ...
func NewDBConfig() (*DBConfig, error) {
	viper.GetViper().SetEnvPrefix("db")

	host := viper.GetString("host")
	if host == "" {
		log.Info("undefined db host, defaulting to localhost")
		host = "localhost"
	}

	port := viper.GetInt("port")
	if port == 0 {
		log.Info("undefined db port, defaulting to 5432")
		port = 5432
	}

	database := viper.GetString("database")
	if database == "" {
		return nil, errors.New("undefined db database")
	}

	user := viper.GetString("user")
	if user == "" {
		return nil, errors.New("undefined db user")
	}

	pass := viper.GetString("password")
	if pass == "" {
		return nil, errors.New("undefined db password")
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		Database: database,
		User:     user,
		Pass:     pass,
	}, nil
}
