package mocks

import (
	"time"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/google/uuid"
)

var mockSnippet = &models.Snippet{
	ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
	Title:     "An old silent pond",
	Content:   "An old silent pond...",
	CreatedOn: time.Now(),
	ExpiresOn: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expireVal int) (string, error) {
	return uuid.New().String(), nil
	// return "9c1fe9ac-b67c-4ba5-9530-208ac6985e0d", nil
}

func (m *SnippetModel) Get(id uuid.UUID) (*models.Snippet, error) {
	switch id {
	case uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"):
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
