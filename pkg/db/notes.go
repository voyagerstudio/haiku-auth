package db

import (
	"errors"
	"fmt"
	"time"
)

// NoteList contains a list of note IDs
type NoteList struct {
	Notes []string `json:"notes"`
}

// Note contains all details needed to dispaly a note
type Note struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetNoteList return a list of note IDs for a given user ID
func (c *Conn) GetNoteList(user string) (*NoteList, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}

	res, err := c.conn.Query("SELECT n.id FROM notes AS n JOIN users AS u ON n.owner_id = u.id WHERE u.id = ?", user)
	if err != nil {
		return nil, fmt.Errorf("error querying for notes: %v", err)
	}
	defer res.Close()

	notes := []string{}
	for res.Next() {
		var noteID string
		if err := res.Scan(&noteID); err != nil {
			return nil, fmt.Errorf("error scanning results: %v", err)
		}
		notes = append(notes, noteID)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error while parsing rows: %v", err)
	}

	return &NoteList{Notes: notes}, nil
}

// GetNote returns a detailed note for a given note ID
func (c *Conn) GetNote(user string, note string) (*Note, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}
	if note == "" {
		return nil, errors.New("note is empty")
	}

	var text string
	var order int
	var createdAt, updatedAt time.Time
	err := c.conn.QueryRow("SELECT data, sort_order, created_at, updated_at FROM notes WHERE id = ?", note).Scan(&text, &order, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning results: %v", err)
	}

	return &Note{
		ID:        note,
		Text:      text,
		Order:     order,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
