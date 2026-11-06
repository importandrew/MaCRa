package player

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestUpdate(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		id          int
		name        string

		expectedErr error
	}{
		"ok": {
			preAddUsers: []int{1},
			id:          1,
			name:        "Alice",
		},
		"not found": {
			id:          1,
			name:        "Alice",
			expectedErr: models.ErrPlayerNotFound,
		},
		// ASCII, so the column rejects it whatever the database encoding is.
		"longer than the column": {
			preAddUsers: []int{1},
			id:          1,
			name:        strings.Repeat("a", 51),
			expectedErr: models.ErrLongName,
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)

			ctx := context.Background()
			conn, err := db.Acquire(ctx)
			if err != nil {
				t.Fatalf("Failed to acquire connection: %v", err)
			}
			defer conn.Release()

			for _, id := range testCase.preAddUsers {
				createTestUser(t, conn, id)
			}

			err = repo.Update(ctx, testCase.id, testCase.name)
			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
