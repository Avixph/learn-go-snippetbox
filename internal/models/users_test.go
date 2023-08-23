package models

import (
	"testing"

	"github.com/Avixph/learn-go-snippetbox/internal/assert"
	"github.com/google/uuid"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	
	// Setup a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID uuid.UUID
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"),
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: uuid.Nil,
			want:   false,
		},
		{
			name:   "Non-existant ID",
			userID: uuid.MustParse("9c1fe9ac-b67c-4ba5-9530-208ac6985e0d"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper func to get a connection pool to the
			// test database. Calling it inside t.Run() means that the fresh database
			// tables and data will be setup and torn down for each sub-test.
			db := newTestDB(t)

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return value and
			// error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}

}
