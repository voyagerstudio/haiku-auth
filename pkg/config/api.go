package config

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type APIConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewAPIConfig
func NewAPIConfig() (*APIConfig, error) {
	viper.GetViper().SetEnvPrefix("api")

	host := viper.GetString("host")
	if host == "" {
		log.Info("undefined api host, defaulting to :http")
	}

	port := viper.GetInt("port")
	if port == 0 {
		log.Info("undefined api port, defaulting to 8080")
		port = 8080
	}

	rt := viper.GetDuration("read_timeout")
	if rt == 0 {
		log.Info("undefined api read timeout, defaulting to 30s")
		rt = 30 * time.Second
	}

	wt := viper.GetDuration("write_timeout")
	if wt == 0 {
		log.Info("undefined api write timeout, defaulting to 30s")
		wt = 30 * time.Second
	}

	return &APIConfig{
		Host:         host,
		Port:         port,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}, nil
}
