package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Add an ErrInvalidCredentials error that returns if a user tries to login
	// with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// Add an ErrDuplicateEmail error that returns if a user tries to signup
	// with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
