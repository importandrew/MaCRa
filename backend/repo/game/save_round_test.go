package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gocognit
func TestSaveRoundOutcome(t *testing.T) {
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
		gameID        int
		roundData     models.RoundData

		expectedErr assert.ErrorAssertionFunc
		twiceErr    assert.ErrorAssertionFunc // if not nil, call SaveRound 2nd time and assert
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID: 1,
			roundData: models.RoundData{
				Round:           1,
				Outcome:         models.Deal,
				WinnerID:        1,
				LoserID:         2,
				Value:           8,
				EastStartScore:  32,
				SouthStartScore: -16,
				WestStartScore:  -8,
				NorthStartScore: -8,
				ForceNextRound:  false,
			},
			expectedErr: assert.NoError,
		},
		"game not exists": {
			gameID:      1,
			expectedErr: assert.Error,
		},
		"deal twice": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID: 1,
			roundData: models.RoundData{
				Round:           1,
				Outcome:         models.Deal,
				WinnerID:        1,
				LoserID:         2,
				Value:           8,
				EastStartScore:  32,
				SouthStartScore: -16,
				WestStartScore:  -8,
				NorthStartScore: -8,
				ForceNextRound:  false,
			},
			expectedErr: assert.NoError,
			twiceErr:    assert.Error,
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

			err = repo.SaveRoundOutcome(ctx, testCase.gameID, testCase.roundData)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.twiceErr != nil {
				err = repo.SaveRoundOutcome(ctx, testCase.gameID, testCase.roundData)
				testCase.twiceErr(t, err)
				if err != nil {
					return
				}
			}

			round := getRound(t, conn, testCase.gameID, testCase.roundData.Round)

			require.NotNil(t, round)
			assert.Equal(t, testCase.roundData.Round, round.RoundNumber)
			assert.Equal(t, string(testCase.roundData.Outcome), round.Outcome)
			require.True(t, round.WinnerID.Valid, "WinnerID should be valid")
			assert.Equal(t, testCase.roundData.WinnerID, int(round.WinnerID.Int64))
			require.True(t, round.LoserID.Valid, "LoserID should be valid")
			assert.Equal(t, testCase.roundData.LoserID, int(round.LoserID.Int64))
			require.True(t, round.Value.Valid, "Value should be valid")
			assert.Equal(t, testCase.roundData.Value, int(round.Value.Int64))
			assert.Equal(t, testCase.roundData.ForceNextRound, round.ForceNextRound)

			scores := getRoundScores(t, conn, round.ID)

			require.Len(t, scores, 4)
			assert.Equal(t, testCase.roundData.EastStartScore, scores["east"])
			assert.Equal(t, testCase.roundData.SouthStartScore, scores["south"])
			assert.Equal(t, testCase.roundData.WestStartScore, scores["west"])
			assert.Equal(t, testCase.roundData.NorthStartScore, scores["north"])
		})
	}
}
