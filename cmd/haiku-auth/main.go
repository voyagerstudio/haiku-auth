package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/voyagerstudio/haiku-auth/pkg/api"
	"github.com/voyagerstudio/haiku-auth/pkg/config"
	"github.com/voyagerstudio/haiku-auth/pkg/db"
)

func main() {
	cfg, err := config.DefaultConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := db.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.Database)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	srv := api.NewServer(cfg.API.Host, cfg.API.Port, db)
	srv.ListenAndServe()
}
