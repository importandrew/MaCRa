package game

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestAbort(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	testData := map[string]struct {
		status            models.GameStatus
		rounds            int
		onlyWithoutRounds bool

		wantErr    error
		wantStatus models.GameStatus
	}{
		"referee aborts a game with rounds": {
			status: models.GameLive, rounds: 2,
			wantStatus: models.GameAborted,
		},
		"seated player aborts a game without rounds": {
			status: models.GameLive, onlyWithoutRounds: true,
			wantStatus: models.GameAborted,
		},
		"seated player can't once a round is in": {
			status: models.GameLive, rounds: 1, onlyWithoutRounds: true,
			wantErr: models.ErrGameHasRounds, wantStatus: models.GameLive,
		},
		"an ended game can't be aborted": {
			status:  models.GameEnded,
			wantErr: models.ErrGameNotLive, wantStatus: models.GameEnded,
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)
			ctx := context.Background()
			conn, err := db.Acquire(ctx)
			require.NoError(t, err)
			defer conn.Release()

			for id := 1; id <= 4; id++ {
				createTestUser(t, conn, id)
			}
			createTestGameWithStatus(t, conn, 1, testCase.status)
			createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
			for n := 1; n <= testCase.rounds; n++ {
				createTestRound(t, conn, 1, preAddRound{roundNumber: n, outcome: models.Draw})
			}

			err = repo.Abort(ctx, 1, testCase.onlyWithoutRounds)
			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.wantStatus, getGameStatus(t, conn, 1))
		})
	}
}

func TestSummary(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "summary", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for id := 1; id <= 4; id++ {
			createTestUser(t, conn, id)
		}
		createTestGame(t, conn, 1)
		createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
		createTestRound(t, conn, 1, preAddRound{roundNumber: 1, outcome: models.Draw})
		createTestRound(t, conn, 1, preAddRound{roundNumber: 2, outcome: models.Draw})

		summary, err := repo.Summary(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, models.GameLive, summary.Status)
		assert.Nil(t, summary.FinishedAt)
		assert.Equal(t, 2, summary.LogSize)

		_, err = repo.Summary(ctx, 999)
		assert.ErrorIs(t, err, models.ErrGameNotFound)
	})
}
