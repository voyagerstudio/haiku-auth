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
	Order     float32   `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteDetailList struct {
	Notes []Note `json:"notes"`
}

// GetNoteList return a list of note IDs for a given user ID
func (c *Conn) GetNoteList(user string) (*NoteList, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}

	res, err := c.conn.Query("SELECT n.id FROM notes AS n JOIN users AS u ON n.owner_id = u.id WHERE u.id = ? ORDER BY n.sort_order", user)
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

// GetNoteList return a list of note IDs for a given user ID
func (c *Conn) GetNoteDetailList(user string) (*NoteDetailList, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}

	res, err := c.conn.Query("SELECT n.id, n.data, n.sort_order, n.created_at, n.updated_at FROM notes AS n JOIN users AS u ON n.owner_id = u.id WHERE u.id = ? ORDER BY n.sort_order", user)
	if err != nil {
		return nil, fmt.Errorf("error querying for notes: %v", err)
	}
	defer res.Close()

	notes := []Note{}
	for res.Next() {
		n := Note{}
		if err := res.Scan(&n.ID, n.Text, n.Order, n.CreatedAt, n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning results: %v", err)
		}
		notes = append(notes, n)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error while parsing rows: %v", err)
	}

	return &NoteDetailList{Notes: notes}, nil
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
	var order float32
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

func (c *Conn) CreateNote(user string, id string, text string, order float32) (*Note, error) {
	if user == "" {
		return nil, errors.New("user is empty")
	}
	if id == "" {
		return nil, errors.New("id is empty")
	}
	if text == "" {
		return nil, errors.New("text is empty")
	}

	var note Note
	err := c.conn.
		QueryRow(
			"INSERT INTO notes SET id, owner_id, data, order = (?, ?, ?, ?) RETURNING data, sort_order, created_at, updated_at",
			id, user, text, order).
		Scan(&note.Text, &note.Order, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error saving note: %v", err)
	}

	return &note, nil
}
