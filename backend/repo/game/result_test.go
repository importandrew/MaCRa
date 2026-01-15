package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameResult(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddPlayers []int
		preAddGame    struct {
			gameID    int
			playerIDs [4]int
		}
		finalData models.FinalGameData

		gameID      int
		expectedErr assert.ErrorAssertionFunc
		checkResult func(*testing.T, models.GameResult)
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGame: struct {
				gameID    int
				playerIDs [4]int
			}{
				gameID:    1,
				playerIDs: [4]int{1, 2, 3, 4},
			},
			finalData: models.FinalGameData{
				GameID: 1,
				East: models.FinalPlayerData{
					ID: 1,
					GP: 200,
					TP: 4,
				},
				South: models.FinalPlayerData{
					ID: 2,
					GP: 50,
					TP: 3.0 / 2.0,
				},
				West: models.FinalPlayerData{
					ID: 3,
					GP: 50,
					TP: 3.0 / 2.0,
				},
				North: models.FinalPlayerData{
					ID: 4,
					GP: -300,
					TP: 0,
				},
			},
			gameID:      1,
			expectedErr: assert.NoError,
			checkResult: func(t *testing.T, result models.GameResult) {
				assert.Equal(t, 1, result.ID)
				require.Len(t, result.PlayerResults, 4)
				// Results should be sorted by TP descending
				assert.Equal(t, 1, result.PlayerResults[0].ID)
				assert.Equal(t, float32(4), result.PlayerResults[0].TP)
				assert.Equal(t, 2, result.PlayerResults[1].ID)
				assert.Equal(t, float32(3.0/2.0), result.PlayerResults[1].TP)
				assert.Equal(t, 3, result.PlayerResults[2].ID)
				assert.Equal(t, float32(3.0/2.0), result.PlayerResults[2].TP)
				assert.Equal(t, 4, result.PlayerResults[3].ID)
				assert.Equal(t, float32(0), result.PlayerResults[3].TP)
			},
		},
		"empty result": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGame: struct {
				gameID    int
				playerIDs [4]int
			}{
				gameID:    1,
				playerIDs: [4]int{1, 2, 3, 4},
			},
			// Don't finalize the game, so no results - returns empty slice, not error
			finalData: models.FinalGameData{},
			gameID:      1,
			expectedErr: assert.NoError,
			checkResult: func(t *testing.T, result models.GameResult) {
				assert.Equal(t, 1, result.ID)
				assert.Empty(t, result.PlayerResults)
			},
		},
		"non-existent game": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGame: struct {
				gameID    int
				playerIDs [4]int
			}{
				gameID:    1,
				playerIDs: [4]int{1, 2, 3, 4},
			},
			// Don't finalize - testing non-existent game
			finalData: models.FinalGameData{},
			gameID:      999,
			expectedErr: assert.NoError,
			checkResult: func(t *testing.T, result models.GameResult) {
				assert.Equal(t, 999, result.ID)
				assert.Empty(t, result.PlayerResults)
			},
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

			createTestGame(t, conn, testCase.preAddGame.gameID)
			createTestGamePlayers(t, conn, testCase.preAddGame.gameID, testCase.preAddGame.playerIDs)

			// Finalize the game if finalData is provided
			if testCase.finalData.GameID > 0 {
				err = repo.Finalize(ctx, testCase.finalData)
				require.NoError(t, err, "Failed to finalize game")
			}

			result, err := repo.GameResult(ctx, testCase.gameID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkResult != nil {
				testCase.checkResult(t, result)
			}
		})
	}
}
