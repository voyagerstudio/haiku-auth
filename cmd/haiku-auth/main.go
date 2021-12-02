package main

import (
	"github.com/sirupsen/logrus"
	"github.com/voyagerstudio/haiku-auth/pkg/api"
	"github.com/voyagerstudio/haiku-auth/pkg/config"
)

func main() {
	cfg, err := config.DefaultConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	srv := api.NewServer(cfg.API.Host, cfg.API.Port)
	srv.ListenAndServe()
}
