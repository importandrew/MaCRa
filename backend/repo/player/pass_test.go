package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestGetPassByEmail(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		email       string

		expectedErr      assert.ErrorAssertionFunc
		expectedPassword string
		expectedErrValue error
	}{
		"ok": {
			preAddUsers:      []int{1, 2},
			email:            "1@test.com",
			expectedErr:      assert.NoError,
			expectedPassword: "password123",
		},
		"not found": {
			preAddUsers:      []int{1},
			email:            "nonexistent@test.com",
			expectedErr:      assert.Error,
			expectedErrValue: models.ErrPlayerNotFound,
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

			password, err := repo.GetPassByEmail(ctx, testCase.email)
			testCase.expectedErr(t, err)
			if err != nil {
				if testCase.expectedErrValue != nil {
					assert.Equal(t, testCase.expectedErrValue, err)
				}
				return
			}

			assert.Equal(t, testCase.expectedPassword, password)
		})
	}
}
