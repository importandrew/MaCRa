package game

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

// TestUndoLastRound_tournament: in a tournament game, a referee may revive an
// ended game however long ago it ended while the tour is live, and nobody
// may undo once the tour is finished.
func TestUndoLastRound_tournament(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tournament undo", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for _, id := range []int{1, 2, 3, 4} {
			createTestUser(t, conn, id)
		}
		createTestGame(t, conn, 1)
		createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})
		createTestRound(t, conn, 1, preAddRound{roundNumber: 1, outcome: models.Draw})
		finalize(t, repo, 1)
		// Ended long ago, in a live tour.
		_, err = conn.Exec(ctx, `
			INSERT INTO tours (id, event_id, number, status, seating, published_at)
			VALUES (7, $1, 1, 'live', 'random', now())`, testEventID)
		require.NoError(t, err)
		_, err = conn.Exec(ctx, `UPDATE games SET tour_id = 7, finished_at = now() - interval '1 day' WHERE id = 1`)
		require.NoError(t, err)

		cutoff := time.Now().Add(-models.UndoFinalRoundTimeout)
		assert.ErrorIs(t, repo.UndoLastRound(ctx, 1, 1, 2, nil), models.ErrGameNotLive, "a seated player doesn't revive")
		require.NoError(t, repo.UndoLastRound(ctx, 1, 1, 1, &cutoff), "a referee revives it while the tour is live")
		assert.Equal(t, models.GameLive, getGameStatus(t, conn, 1))

		_, err = conn.Exec(ctx, `UPDATE tours SET status = 'finished' WHERE id = 7`)
		require.NoError(t, err)
		createTestRound(t, conn, 1, preAddRound{roundNumber: 1, outcome: models.Draw})
		assert.ErrorIs(t, repo.UndoLastRound(ctx, 1, 1, 1, &cutoff), models.ErrTourFinished)
	})
}
