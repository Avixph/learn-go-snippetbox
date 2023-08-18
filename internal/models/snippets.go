package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Define a SnippetModelInterface interface that describes the methods our
// SnippetModel has.
type SnippetModelInterface interface {
	Insert(title string, content string, expireVal int) (string, error)
	Get(id uuid.UUID) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// Define a Snippet type that holds data for individual snippets. Notice
// how the feilds of the struct correspond to the feilds in our PostgreSQL
// snippets table?
type Snippet struct {
	ID        uuid.UUID
	Title     string
	Content   string
	CreatedOn time.Time
	ExpiresOn time.Time
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

	// The id returned has the type uuid, so we convert it to a string type
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
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedOn, &s.ExpiresOn)
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
	// Define the SQL query we want to execute.
	query := `SELECT id, title, content, created_on, expires_on FROM snippets
	WHERE expires_on > now() ORDER BY id LIMIT 10`

	// Use the Query() method on the connection pool to execute the query.
	// This returns a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is alays
	// properly closed before the Latest() method returns. This defer
	// statement should come *after* checking for an err from the Query()
	// method. Otherwise, if Query() returns an err, we'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and subsequent) row to be acted on by the rows.
	// Scan() method. If iteration over all the rows completes then the
	// resultset autoatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the ow to the
		// new Snippet object that we created. Agian, the arguments to row.Scan
		// () must be pointers to the place we want to copy the data into, and
		// the number of arguments must be exactly the same as the number of
		// columns returned by the statement.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedOn, &s.ExpiresOn)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve
	// any error that was encounterd during the iteration. It's important to
	// call this - don't asume that a successful iteration was completd over
	// the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went ok, then return the Snippets slice.
	return snippets, nil
}
