package game

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
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
			twiceErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, models.ErrLogChanged)
			},
		},
		"penalty twice from the same log": {
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
				Outcome:         models.Penalty,
				LoserID:         2,
				Value:           10,
				EastStartScore:  0,
				SouthStartScore: -10,
				LogSize:         0,
			},
			expectedErr: assert.NoError,
			twiceErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, models.ErrLogChanged)
			},
		},
		"stale log size": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID: 1,
			roundData: models.RoundData{
				Round:   1,
				Outcome: models.Draw,
				LogSize: 1,
			},
			expectedErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, models.ErrLogChanged)
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

// TestSaveRoundOutcome_ConcurrentPenalties checks that two entries computed
// from the same log (a double tap on two phones) can't both be saved.
func TestSaveRoundOutcome_ConcurrentPenalties(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.Test(t, "concurrent penalties", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()

		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		for _, id := range []int{1, 2, 3, 4} {
			createTestUser(t, conn, id)
		}
		createTestGame(t, conn, 1)
		createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
		// Free this connection: each attempt takes its own.
		conn.Release()

		data := models.RoundData{Round: 1, Outcome: models.Penalty, LoserID: 2, Value: 10, SouthStartScore: -10}

		const attempts = 5
		errs := make(chan error, attempts)
		var wg sync.WaitGroup
		for range attempts {
			wg.Go(func() { errs <- repo.SaveRoundOutcome(ctx, 1, data) })
		}
		wg.Wait()
		close(errs)

		saved := 0
		for err := range errs {
			if err == nil {
				saved++
				continue
			}
			assert.ErrorIs(t, err, models.ErrLogChanged)
		}
		assert.Equal(t, 1, saved)
	})
}
