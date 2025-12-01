package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameLog(t *testing.T) {
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
		preAddRounds  []preAddRound
		gameID        int

		expectedErr  assert.ErrorAssertionFunc
		expectedLogs []models.RoundLog
	}{
		"ok empty": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID:       1,
			expectedErr:  assert.NoError,
			expectedLogs: []models.RoundLog{},
		},
		"ok many rounds": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			preAddRounds: []preAddRound{
				{
					roundNumber: 1,
					outcome:     models.Deal,
					winnerID:    1,
					loserID:     2,
					value:       8,
					scores: map[string]int{
						"east":  32,
						"south": -16,
						"west":  -8,
						"north": -8,
					},
				},
				{
					roundNumber: 2,
					outcome:     models.Self,
					winnerID:    3,
					loserID:     0,
					value:       8,
					scores: map[string]int{
						"east":  24,
						"south": -8,
						"west":  40,
						"north": -8,
					},
				},
				{
					roundNumber: 3,
					outcome:     models.Draw,
					winnerID:    0,
					loserID:     0,
					value:       0,
					scores: map[string]int{
						"east":  24,
						"south": -8,
						"west":  40,
						"north": -8,
					},
				},
			},
			gameID:      1,
			expectedErr: assert.NoError,
			expectedLogs: []models.RoundLog{
				{
					Round:          1,
					EventSequence:  1,
					Outcome:        models.Deal,
					Winner:         1,
					Loser:          2,
					Value:          8,
					ForceNextRound: false,
					EastScore:      32,
					SouthScore:     -16,
					WestScore:      -8,
					NorthScore:     -8,
				},
				{
					Round:          2,
					EventSequence:  1,
					Outcome:        models.Self,
					Winner:         3,
					Loser:          0,
					Value:          8,
					ForceNextRound: false,
					EastScore:      24,
					SouthScore:     -8,
					WestScore:      40,
					NorthScore:     -8,
				},
				{
					Round:          3,
					EventSequence:  1,
					Outcome:        models.Draw,
					Winner:         0,
					Loser:          0,
					Value:          0,
					ForceNextRound: false,
					EastScore:      24,
					SouthScore:     -8,
					WestScore:      40,
					NorthScore:     -8,
				},
			},
		},
		"ok with penalty": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			preAddRounds: []preAddRound{
				{
					roundNumber: 2,
					outcome:     models.Self,
					winnerID:    1,
					value:       12,
					scores: map[string]int{
						"east":  44,
						"south": -28,
						"west":  -38,
						"north": 12,
					},
				},
				{
					roundNumber: 1,
					outcome:     models.Penalty,
					loserID:     3,
					value:       10,
					scores: map[string]int{
						"east":  0,
						"south": 0,
						"west":  -10,
						"north": 0,
					},
				},
				{
					roundNumber: 1,
					outcome:     models.Deal,
					winnerID:    4,
					loserID:     1,
					value:       8,
					scores: map[string]int{
						"east":  -16,
						"south": -8,
						"west":  -18,
						"north": 32,
					},
				},
			},
			gameID:      1,
			expectedErr: assert.NoError,
			expectedLogs: []models.RoundLog{
				{
					Round:          1,
					EventSequence:  1,
					Outcome:        models.Penalty,
					Winner:         0,
					Loser:          3,
					Value:          10,
					ForceNextRound: false,
					EastScore:      0,
					SouthScore:     0,
					WestScore:      -10,
					NorthScore:     0,
				},
				{
					Round:          1,
					EventSequence:  2,
					Outcome:        models.Deal,
					Winner:         4,
					Loser:          1,
					Value:          8,
					ForceNextRound: false,
					EastScore:      -16,
					SouthScore:     -8,
					WestScore:      -18,
					NorthScore:     32,
				},
				{
					Round:          2,
					EventSequence:  1,
					Outcome:        models.Self,
					Winner:         1,
					Value:          12,
					ForceNextRound: false,
					EastScore:      44,
					SouthScore:     -28,
					WestScore:      -38,
					NorthScore:     12,
				},
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

			for _, cfg := range testCase.preAddGames {
				createTestGame(t, conn, cfg.gameID)
				createTestGamePlayers(t, conn, cfg.gameID, cfg.playerIDs)
			}

			for _, round := range testCase.preAddRounds {
				createTestRound(t, conn, testCase.gameID, round)
			}

			logs, err := repo.GameLog(ctx, testCase.gameID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Len(t, logs, len(testCase.expectedLogs), "Number of logs should match")

			for i, expectedLog := range testCase.expectedLogs {
				actualLog := logs[i]
				assert.Positive(t, actualLog.ID, "Log ID should be greater than 0")
				assert.Equal(t, expectedLog.Round, actualLog.Round, "Round number should match")
				assert.Equal(t, expectedLog.EventSequence, actualLog.EventSequence, "Event sequence should match")
				assert.Equal(t, expectedLog.Outcome, actualLog.Outcome, "Outcome should match")
				assert.Equal(t, expectedLog.Winner, actualLog.Winner, "Winner should match")
				assert.Equal(t, expectedLog.Loser, actualLog.Loser, "Loser should match")
				assert.Equal(t, expectedLog.Value, actualLog.Value, "Value should match")
				assert.Equal(t, expectedLog.ForceNextRound, actualLog.ForceNextRound, "ForceNextRound should match")
				assert.Equal(t, expectedLog.EastScore, actualLog.EastScore, "EastScore should match")
				assert.Equal(t, expectedLog.SouthScore, actualLog.SouthScore, "SouthScore should match")
				assert.Equal(t, expectedLog.WestScore, actualLog.WestScore, "WestScore should match")
				assert.Equal(t, expectedLog.NorthScore, actualLog.NorthScore, "NorthScore should match")
				assert.NotZero(t, actualLog.CreatedAt, "CreatedAt should be set")
			}
		})
	}
}
