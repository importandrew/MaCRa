package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		seating       [4]int
		preAddPlayers []int

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers map[string]int // wind -> player_id
	}{
		"ok": {
			seating:       [4]int{1, 2, 3, 4},
			preAddPlayers: []int{1, 2, 3, 4},
			expectedErr:   assert.NoError,
			expectedPlayers: map[string]int{
				"east":  1,
				"south": 2,
				"west":  3,
				"north": 4,
			},
		},
		"player not exists": {
			seating:       [4]int{1, 2, 3, 4},
			preAddPlayers: []int{1, 2, 3, 5},
			expectedErr:   assert.Error,
		},
		"duplicate player": {
			seating:       [4]int{1, 2, 3, 3},
			preAddPlayers: []int{1, 2, 3},
			expectedErr:   assert.Error,
		},
		"short seating": {
			seating:       [4]int{1, 2, 3},
			preAddPlayers: []int{1, 2, 3},
			expectedErr:   assert.Error,
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

			for _, id := range testCase.preAddPlayers {
				createTestUser(t, conn, id)
			}

			gameID, err := repo.NewGame(ctx, testCase.seating)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Positive(t, gameID, "Game ID should be greater than 0")

			status := getGameStatus(t, conn, gameID)
			assert.Equal(t, models.GameLive, status, "Game status should be %s", models.GameLive)

			players := getGamePlayers(t, conn, gameID)
			require.Len(t, players, 4, "Game should have exactly 4 players")
			for wind, id := range testCase.expectedPlayers {
				assert.Equal(t, players[wind], id)
			}
		})
	}
}
