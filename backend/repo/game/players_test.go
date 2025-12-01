package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestPlayers(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type preAddGame struct {
		gameID        int
		playerIDs     [4]int
		gamePlayerIDs [4]int
	}

	testData := map[string]struct {
		preAddGames   []preAddGame
		preAddPlayers []int
		gameID        int

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers models.StartingPlayers
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:        1,
					playerIDs:     [4]int{1, 2, 3, 4},
					gamePlayerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID:      1,
			expectedErr: assert.NoError,
			expectedPlayers: models.StartingPlayers{
				East:  1,
				South: 2,
				West:  3,
				North: 4,
			},
		},
		"not exists": {
			gameID:      1,
			expectedErr: assert.Error,
		},
		"wrong number of records": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:        1,
					playerIDs:     [4]int{1, 2, 3, 4},
					gamePlayerIDs: [4]int{1, 2, 0, 4},
				},
			},
			gameID:      1,
			expectedErr: assert.Error,
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

			for _, cfg := range testCase.preAddGames {
				createTestGame(t, conn, cfg.gameID)
				createTestGamePlayers(t, conn, cfg.gameID, cfg.gamePlayerIDs)
			}

			players, err := repo.Players(ctx, testCase.gameID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedPlayers.East, players.East, "East player should match")
			assert.Equal(t, testCase.expectedPlayers.South, players.South, "South player should match")
			assert.Equal(t, testCase.expectedPlayers.West, players.West, "West player should match")
			assert.Equal(t, testCase.expectedPlayers.North, players.North, "North player should match")
		})
	}
}
