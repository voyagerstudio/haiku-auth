package db

import (
	"errors"
	"fmt"
	"time"
)

// UserList contains a list of user IDs
type UserList struct {
	Users []string `json:"users"`
}

// User contains all details describing a user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserList return a list of user IDs
func (c *Conn) GetUserList() (*UserList, error) {

	res, err := c.conn.Query("SELECT n.id FROM users AS n")
	if err != nil {
		return nil, fmt.Errorf("error querying for users: %v", err)
	}
	defer res.Close()

	users := []string{}
	for res.Next() {
		var userID string
		if err := res.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning results: %v", err)
		}
		users = append(users, userID)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error while parsing rows: %v", err)
	}

	return &UserList{Users: users}, nil
}

// GetUser returns a detailed user for a given note ID
func (c *Conn) GetUser(user string) (*User, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}

	var id string
	var username string
	var email string
	var createdAt, updatedAt time.Time
	err := c.conn.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", user).Scan(&id, &username, &email, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning results: %v", err)
	}

	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (c *Conn) CreateUser(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	var user User
	err := c.conn.
		QueryRow(
			"INSERT INTO users SET id VALUES (?) RETURNING id, created_at, updated_at",
			id).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error saving user: %v", err)
	}

	return &user, nil
}
