package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Define a Snippet type that holds data for individual snippets. Notice
// how the feilds of the struct correspond to the feilds in our PostgreSQL
// snippets table?
type Snippet struct {
	ID         uuid.UUID
	Title      string
	Content    string
	Created_On time.Time
	Expires_On time.Time
}

// Define a SnippetModel type that wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// The Insert() method will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expireVal int) (uuid.UUID, error) {
	// Define the SQL query we want to execute.
	query := `INSERT INTO snippets (title, content, created_on, expires_on)
		VALUES ($1, $2, (now() at time zone 'utc'), (now() at time zone 'utc' + interval '$3 day'))
		RETURNING ID`

	// Create an args slice containing the values for the placeholder
	// parameters. The first parameter is the stmt var, followed by the
	// title, content and the expiry values for the palceholder parameters.
	// Declaring this slice next to our SQL query helps to make it nice and
	// clear *what values are being used where* in the query.
	args := []any{title, content, expireVal}

	// Create an id var with the type uuid.UUID
	var id uuid.UUID

	// Use the QueryRow() method to execute the SQL query on our connection
	// pool, passing in args as a variadic parameter and scanning the
	// generated id.
	row := m.DB.QueryRow(query, args...)
	err := row.Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	// The returned id has a type UUID
	// return id, nil
	return id, nil
}

// The Get() method will return a specific snippet from the database.
func (m *SnippetModel) Get(id uuid.UUID) (*Snippet, error) {
	return nil, nil
}

// The Latest() method will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
