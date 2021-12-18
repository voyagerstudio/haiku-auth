package db

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Conn struct {
	conn *pg.DB
}

// New initializes a new database connection
func New(host string, port int, user, password, database string) (*Conn, error) {
	conn := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		User:     user,
		Password: password,
		Database: database,
	})

	err := conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error pinging db: %v", err)
	}

	return &Conn{
		conn: conn,
	}, nil
}
