package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		playerID    int

		expectedErr      assert.ErrorAssertionFunc
		expectedPlayer   models.Player
		expectedErrValue error
	}{
		"ok": {
			preAddUsers:    []int{1, 2},
			playerID:       1,
			expectedErr:    assert.NoError,
			expectedPlayer: models.Player{ID: 1, Name: "1"},
		},
		"not found": {
			preAddUsers:      []int{1},
			playerID:         2,
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

			player, err := repo.Get(ctx, testCase.playerID)
			testCase.expectedErr(t, err)
			if err != nil {
				if testCase.expectedErrValue != nil {
					assert.Equal(t, testCase.expectedErrValue, err)
				}
				return
			}

			assert.Equal(t, testCase.expectedPlayer.ID, player.ID)
			assert.Equal(t, testCase.expectedPlayer.Name, player.Name)
		})
	}
}

func TestGetPlayerByEmail(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		email       string

		expectedErr      assert.ErrorAssertionFunc
		expectedPlayer   models.Player
		expectedErrValue error
	}{
		"ok": {
			preAddUsers:    []int{1, 2},
			email:          "1@test.com",
			expectedErr:    assert.NoError,
			expectedPlayer: models.Player{ID: 1, Name: "1"},
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

			player, err := repo.GetByEmail(ctx, testCase.email)
			testCase.expectedErr(t, err)
			if err != nil {
				if testCase.expectedErrValue != nil {
					assert.Equal(t, testCase.expectedErrValue, err)
				}
				return
			}

			assert.Equal(t, testCase.expectedPlayer.ID, player.ID)
			assert.Equal(t, testCase.expectedPlayer.Name, player.Name)
		})
	}
}
