package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created_on)
		VALUES ($1, $2, $3, (now() at time zone 'utc'))
		RETURNING ID`

	args := []any{name, email, hashedPassword}

	var id uuid.UUID

	// Use the QueryRow() method to insert the user details and hashed password into the user table.
	row := m.DB.QueryRow(query, args...)
	err = row.Scan(&id)
	if err != nil {
		// If this returns an err, we use the errors.As() func to check
		// whether the error has the type *pq.ErrprCpde. If it does,
		// the error will be assigned to the pSQLError variable. We
		// can then check whether or not the error relates to our
		// users_uc_email key by checking if the error code equals and
		// the contents of the err mesage string. IF it does, we return
		// an ErrDuplicateEmail error.
		var pSQLError *pq.Error
		if errors.As(err, &pSQLError) {
			if pSQLError.Code == "23505" && strings.Contains(pSQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// The Authenticate() method will verify whether a user with the
// provided email and password exists. If they do the relevent user
// ID will be returned.
func (m *UserModel) Authenticate(email, password string) (string, error) {
	// Retrieve the id and hashed password associated withthe given email.
	// If no matching email exists then we return the ErrInvalidCredentials
	// error.
	var id uuid.UUID
	var hashedPassword []byte

	query := `SELECT id, hashed_password FROM users WHERE email = $1`

	row := m.DB.QueryRow(query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil.String(), ErrInvalidCredentials
		} else {
			return uuid.Nil.String(), err
		}
	}

	// Check whether the hashed password and plain-text password match. If
	// they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil.String(), ErrInvalidCredentials
		} else {
			return uuid.Nil.String(), err
		}
	}

	// Else, the password is correct and return the ID
	return id.String(), nil
}

// The Exists() method will checkif a user with a specific ID
// exists.
func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	return false, nil
}
