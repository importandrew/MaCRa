package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

// TestList checks the caller's roles, the last game (aborted games don't
// count) and the order: active events by their last game or creation, then
// finished ones.
func TestList(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "list", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		old := createTestEvent(t, db)      // played an hour ago
		quiet := createTestEvent(t, db)    // no games, created now
		finished := createTestEvent(t, db) // finished
		aborted := createTestEvent(t, db)  // only an aborted game

		gameID, err := game.NewRepo(db).NewGame(ctx, [4]int{1, 2, 3, 4}, old)
		require.NoError(t, err)
		require.NoError(t, game.NewRepo(db).SetStatus(ctx, gameID, models.GameEnded))
		_, err = conn.Exec(ctx, `UPDATE games SET created_at = $1 WHERE id = $2`, time.Now().UTC().Add(-time.Hour), gameID)
		require.NoError(t, err)
		_, err = conn.Exec(ctx, `UPDATE events SET created_at = $1 WHERE id IN ($2, $3)`, time.Now().UTC().Add(-48*time.Hour), old, aborted)
		require.NoError(t, err)
		require.NoError(t, repo.Finish(ctx, finished))
		abortedGame, err := game.NewRepo(db).NewGame(ctx, [4]int{1, 2, 3, 4}, aborted)
		require.NoError(t, err)
		require.NoError(t, game.NewRepo(db).Abort(ctx, abortedGame, false))

		// User 1 owns the test events: player and referee.
		rows, err := repo.List(ctx, 1)
		require.NoError(t, err)

		ids := make([]int, len(rows))
		for i, row := range rows {
			ids[i] = row.ID
		}
		assert.Equal(t, []int{quiet, old, aborted, finished}, ids)

		byID := map[int]models.EventListRow{}
		for _, row := range rows {
			byID[row.ID] = row
		}
		assert.NotNil(t, byID[old].LastGameAt)
		assert.Nil(t, byID[quiet].LastGameAt)
		assert.Nil(t, byID[aborted].LastGameAt, "an aborted game doesn't count")
		assert.Equal(t, []models.EventRole{models.EventPlayer, models.EventReferee}, byID[old].Roles)

		visitor, err := repo.List(ctx, 0)
		require.NoError(t, err)
		assert.Empty(t, visitor[0].Roles)
	})
}
