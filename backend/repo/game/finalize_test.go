package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinalize(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type preAddGame struct {
		gameID    int
		playerIDs [4]int
	}

	testData := map[string]struct {
		preAddGames   []preAddGame
		preAddPlayers []int
		finalData     models.FinalGameData

		expectedErr assert.ErrorAssertionFunc
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
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
					TP: float32(3) / 2,
				},
				West: models.FinalPlayerData{
					ID: 3,
					GP: 50,
					TP: float32(3) / 2,
				},
				North: models.FinalPlayerData{
					ID: 4,
					GP: -300,
					TP: 0,
				},
			},
			expectedErr: assert.NoError,
		},
		"game not exists": {
			finalData: models.FinalGameData{
				GameID: 999,
			},
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
				createTestGamePlayers(t, conn, cfg.gameID, cfg.playerIDs)
			}

			err = repo.Finalize(ctx, testCase.finalData)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			// Verify game status
			status := getGameStatus(t, conn, testCase.finalData.GameID)
			assert.Equal(t, models.GameEnded, status, "Game status should be ended")

			// Verify player results
			results := getPlayerResults(t, conn, testCase.finalData.GameID)
			require.Len(t, results, 4, "Should have 4 player results")

			// Create a map for easier lookup
			resultMap := make(map[int]playerResult)
			for _, r := range results {
				resultMap[r.ID] = r
			}

			// Verify each player's results
			assert.Equal(
				t,
				testCase.finalData.East.GP,
				resultMap[testCase.finalData.East.ID].GP,
				"East GP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.East.TP,
				resultMap[testCase.finalData.East.ID].TP,
				"East TP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.South.GP,
				resultMap[testCase.finalData.South.ID].GP,
				"South GP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.South.TP,
				resultMap[testCase.finalData.South.ID].TP,
				"South TP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.West.GP,
				resultMap[testCase.finalData.West.ID].GP,
				"West GP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.West.TP,
				resultMap[testCase.finalData.West.ID].TP,
				"West TP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.North.GP,
				resultMap[testCase.finalData.North.ID].GP,
				"North GP should match",
			)
			assert.Equal(
				t,
				testCase.finalData.North.TP,
				resultMap[testCase.finalData.North.ID].TP,
				"North TP should match",
			)
		})
	}
}
