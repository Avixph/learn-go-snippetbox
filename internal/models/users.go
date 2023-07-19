package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Define a User type.
type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword []byte
	CreatedOn      time.Time
}

// Define a UserModel type that wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// The Insert() method will add a new recod to the "users" table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// The Authenticate() method will verify whether a user with the
// provided email and password exists. If they do the relevent user
// ID will be returned.
func (m *UserModel) Authenticate(email, password string) (string, error) {
	return uuid.Nil.String(), nil
}

// The Exists() method will checkif a user with a specific ID
// exists.
func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	return false, nil
}
