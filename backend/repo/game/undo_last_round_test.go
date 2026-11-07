package game

import (
	"context"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

//nolint:gocognit
func TestUndoLastRound(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	draw := func(n int) preAddRound { return preAddRound{roundNumber: n, outcome: models.Draw} }

	testData := map[string]struct {
		rounds  int
		setup   func(t *testing.T, repo *Repo, conn *pgxpool.Conn)
		logSize int

		expectedErr     error
		expectedErrText string
		expectedStatus  models.GameStatus
		expectedRounds  int
	}{
		"live game": {
			rounds:         3,
			logSize:        3,
			expectedStatus: models.GameLive,
			expectedRounds: 2,
		},
		"log changed since it was loaded": {
			rounds:         3,
			logSize:        2,
			expectedErr:    models.ErrLogChanged,
			expectedStatus: models.GameLive,
			expectedRounds: 3,
		},
		"no rounds": {
			rounds:         0,
			logSize:        0,
			expectedErr:    models.ErrNoRounds,
			expectedStatus: models.GameLive,
		},
		"ended game becomes live again": {
			rounds:  16,
			logSize: 16,
			setup: func(t *testing.T, repo *Repo, _ *pgxpool.Conn) {
				finalize(t, repo, 1)
			},
			expectedStatus: models.GameLive,
			expectedRounds: 15,
		},
		"ended game too late": {
			rounds:  16,
			logSize: 16,
			setup: func(t *testing.T, repo *Repo, conn *pgxpool.Conn) {
				finalize(t, repo, 1)
				_, err := conn.Exec(context.Background(),
					`UPDATE games SET finished_at = now() - interval '2 hours' WHERE id = 1`)
				require.NoError(t, err)
			},
			expectedErr:    models.ErrTooLate,
			expectedStatus: models.GameEnded,
			expectedRounds: 16,
		},
		"ended game, a player has started another game": {
			rounds:  16,
			logSize: 16,
			setup: func(t *testing.T, repo *Repo, conn *pgxpool.Conn) {
				finalize(t, repo, 1)
				createTestUser(t, conn, 5)
				createTestGame(t, conn, 2)
				createTestGamePlayers(t, conn, 2, [4]int{5, 2, 6, 7})
			},
			expectedErr:     models.ErrUndoFinalRound,
			expectedErrText: "2 has already started another game",
			expectedStatus:  models.GameEnded,
			expectedRounds:  16,
		},
		"aborted game": {
			rounds:  3,
			logSize: 3,
			setup: func(t *testing.T, repo *Repo, _ *pgxpool.Conn) {
				require.NoError(t, repo.SetStatus(context.Background(), 1, models.GameAborted))
			},
			expectedErr:    models.ErrGameNotLive,
			expectedStatus: models.GameAborted,
			expectedRounds: 3,
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)
			ctx := context.Background()

			conn, err := db.Acquire(ctx)
			require.NoError(t, err)
			defer conn.Release()

			for _, id := range []int{1, 2, 3, 4, 6, 7} {
				createTestUser(t, conn, id)
			}
			createTestGame(t, conn, 1)
			createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
			for n := 1; n <= testCase.rounds; n++ {
				createTestRound(t, conn, 1, draw(n))
			}
			if testCase.setup != nil {
				testCase.setup(t, repo, conn)
			}

			err = repo.UndoLastRound(ctx, 1, testCase.logSize)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				if testCase.expectedErrText != "" {
					assert.EqualError(t, err, testCase.expectedErrText)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedStatus, getGameStatus(t, conn, 1))
			assert.Equal(t, testCase.expectedRounds, countRounds(t, conn, 1))
			if testCase.expectedErr == nil && testCase.expectedStatus == models.GameLive {
				assert.Empty(t, getPlayerResults(t, conn, 1), "a live game has no results")
			}
		})
	}
}

// TestUndoLastRound_Concurrent checks that two referees pressing Undo at once
// remove one entry, not two.
func TestUndoLastRound_Concurrent(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.Test(t, "concurrent undo", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()

		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		for _, id := range []int{1, 2, 3, 4} {
			createTestUser(t, conn, id)
		}
		createTestGame(t, conn, 1)
		createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
		for n := 1; n <= 3; n++ {
			createTestRound(t, conn, 1, preAddRound{roundNumber: n, outcome: models.Draw})
		}
		// Free this connection: each attempt takes its own.
		conn.Release()

		const attempts = 5
		errs := make(chan error, attempts)
		var wg sync.WaitGroup
		for range attempts {
			wg.Go(func() { errs <- repo.UndoLastRound(ctx, 1, 3) })
		}
		wg.Wait()
		close(errs)

		undone := 0
		for err := range errs {
			if err == nil {
				undone++
				continue
			}
			assert.ErrorIs(t, err, models.ErrLogChanged)
		}
		assert.Equal(t, 1, undone)

		check, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer check.Release()
		assert.Equal(t, 2, countRounds(t, check, 1))
	})
}

func finalize(t *testing.T, repo *Repo, gameID int) {
	t.Helper()
	err := repo.Finalize(context.Background(), models.FinalGameData{
		GameID: gameID,
		East:   models.FinalPlayerData{ID: 1, TP: 1},
		South:  models.FinalPlayerData{ID: 2, TP: 1},
		West:   models.FinalPlayerData{ID: 3, TP: 1},
		North:  models.FinalPlayerData{ID: 4, TP: 1},
	})
	require.NoError(t, err)
}

func countRounds(t *testing.T, conn *pgxpool.Conn, gameID int) int {
	t.Helper()
	rows, err := conn.Query(context.Background(), `SELECT count(*) FROM rounds WHERE game_id = $1`, gameID)
	require.NoError(t, err)
	n, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	require.NoError(t, err)
	return n
}
