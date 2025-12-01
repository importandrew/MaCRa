package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers   []int
		emailToDelete string

		expectedErr      assert.ErrorAssertionFunc
		expectedErrValue error
	}{
		"ok": {
			preAddUsers:   []int{1, 2},
			emailToDelete: "1@test.com",
			expectedErr:   assert.NoError,
		},
		"not found": {
			preAddUsers:      []int{1},
			emailToDelete:    "nonexistent@test.com",
			expectedErr:      assert.Error,
			expectedErrValue: ErrPlayerNotFound,
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

			err = repo.Delete(ctx, testCase.emailToDelete)
			testCase.expectedErr(t, err)
			if err != nil {
				if testCase.expectedErrValue != nil {
					assert.Equal(t, testCase.expectedErrValue, err)
				}
				return
			}

			// Verify the user was actually deleted
			var count int
			err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", testCase.emailToDelete).Scan(&count)
			assert.NoError(t, err, "Should be able to query users")
			assert.Equal(t, 0, count, "User should be deleted")
		})
	}
}
