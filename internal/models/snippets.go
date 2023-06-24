package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Define a Snippet type that holds data for individual snippets. Notice how
// the feilds of the struct correspond to the feilds in our PostgreSQL
// snippets table?
type Snippet struct {
	ID         uuid.UUID
	Title      string
	Content    string
	Created_On time.Time
	Updated_On time.Time
	Expires_On time.Time
}

// Define a SnippetModel type that wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// The Insert() method will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// The Get() method will return a specific snippet from the database.
func (m *SnippetModel) Get(id uuid.UUID) (*Snippet, error) {
	return nil, nil
}

// The Latest() method will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
