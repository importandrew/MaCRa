package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveGames(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type preAddGame struct {
		gameID    int
		playerIDs [4]int
		status    models.GameStatus
	}

	testData := map[string]struct {
		preAddGames   []preAddGame
		preAddPlayers []int

		expectedErr assert.ErrorAssertionFunc
		expectedIDs []int
	}{
		"ok single": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameLive,
				},
			},
			expectedErr: assert.NoError,
			expectedIDs: []int{1},
		},
		"ok many": {
			preAddPlayers: []int{1, 2, 3, 4, 5, 6, 7, 8},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameLive,
				},
				{
					gameID:    2,
					playerIDs: [4]int{5, 6, 7, 8},
					status:    models.GameLive,
				},
			},
			expectedErr: assert.NoError,
			expectedIDs: []int{1, 2},
		},
		"ok filter": {
			preAddPlayers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameLive,
				},
				{
					gameID:    2,
					playerIDs: [4]int{5, 6, 7, 8},
					status:    models.GameAborted,
				},
				{
					gameID:    3,
					playerIDs: [4]int{9, 10, 11, 12},
					status:    models.GameLive,
				},
			},
			expectedErr: assert.NoError,
			expectedIDs: []int{1, 3},
		},
		"ok zero": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameAborted,
				},
			},
			expectedErr: assert.NoError,
			expectedIDs: []int{},
		},
		"ok empty": {
			preAddPlayers: []int{},
			preAddGames:   []preAddGame{},
			expectedErr:   assert.NoError,
			expectedIDs:   []int{},
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
				createTestGameWithStatus(t, conn, cfg.gameID, cfg.status)
				createTestGamePlayers(t, conn, cfg.gameID, cfg.playerIDs)
			}

			gameIDs, err := repo.LiveGames(ctx)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Len(
				t,
				gameIDs,
				len(testCase.expectedIDs),
				"Expected %d games, got %d",
				len(testCase.expectedIDs),
				len(gameIDs),
			)
			assert.ElementsMatch(t, testCase.expectedIDs, gameIDs, "Game IDs should match")
		})
	}
}
