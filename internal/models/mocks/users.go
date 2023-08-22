package mocks

import (
	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/google/uuid"
)

type UserModel struct{}

var uid = uuid.New()

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "kopi@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	if email == "falso@example.com" && password == "pa$$w0rd8923" {
		return uid.String(), nil
		// return "6ba7b811-9dad-11d1-80b4-00c04fd430c8", nil
	}

	return uuid.Nil.String(), models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	switch id {
	case uid:
		// case uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"):
		return true, nil
	default:
		return false, nil
	}
}