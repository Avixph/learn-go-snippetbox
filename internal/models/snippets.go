package models

import (
	"database/sql"
	"errors"
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
func (m *SnippetModel) Insert(title string, content string, expireVal int) (string, error) {
	// Define the SQL query we want to execute.
	query := `INSERT INTO snippets (title, content, created_on, expires_on)
		VALUES ($1, $2, (now() at time zone 'utc'), (now() at time zone 'utc' + $3 * interval '1 day'))
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
		return uuid.Nil.String(), err
	}

	// The id returned has the tpe uuid, so we convert it to a string type
	// before returning
	return id.String(), nil
}

// The Get() method will return a specific snippet from the database.
func (m *SnippetModel) Get(id uuid.UUID) (*Snippet, error) {
	// Define the SQL query we want to execute.
	query := `SELECT id, title, content, created_on, expires_on FROM snippets
	WHERE expires_on > now() AND id = $1`

	// Use the QueryRow() method on the connection pool to execute the query,
	// passing in the untrusted id variable as the value for the placeholder
	// parameter. This returns a pointer to a sql.Row object which holds the
	// result from the database.
	row := m.DB.QueryRow(query, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan() are *pointers* to the place we want to copy the data
	// into, and the number of arguments must be exactly the same as the
	// number of columns returned by the statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created_On, &s.Expires_On)
	if err != nil {
		// If the query returns no rows, the row.Scan() will return a
		// sql.ErrNoRows err. We use the errors.Is() func to check for that
		// err specificallly, and retun our own ErrNoRecord err instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything goes ok then we return the Snippet object.
	return s, nil
}

// The Latest() method will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
