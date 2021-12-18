package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Conn struct {
	conn *sql.DB
}

// New initializes a new database connection
func New(host string, port int, user, password, database string) (*Conn, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging db: %v", err)
	}

	return &Conn{
		conn: conn,
	}, nil
}
