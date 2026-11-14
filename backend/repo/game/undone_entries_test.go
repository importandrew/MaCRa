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

// TestUndoneEntries checks the audit trail: an entry keeps who entered it,
// and an undo keeps the entry with who undid it.
func TestUndoneEntries(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "undone entries", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for _, id := range []int{1, 2, 3, 4, 6} {
			createTestUser(t, conn, id)
		}
		createTestGame(t, conn, 1)
		createTestGamePlayers(t, conn, 1, [4]int{1, 2, 3, 4})

		undone, err := repo.UndoneEntries(ctx, 1)
		require.NoError(t, err)
		assert.Empty(t, undone)

		before := time.Now().UTC().Add(-time.Second)
		require.NoError(t, repo.SaveRoundOutcome(ctx, 1, models.RoundData{
			Round: 1, Outcome: models.Self, WinnerID: 2, Value: 8,
			EastStartScore: -16, SouthStartScore: 48, WestStartScore: -16, NorthStartScore: -16,
			EnteredBy: 2,
		}))
		// An entry from before authors were recorded has none.
		createTestRound(t, conn, 1, preAddRound{roundNumber: 2, outcome: models.Draw,
			scores: map[string]int{"east": -16, "south": 48, "west": -16, "north": -16}})

		log, err := repo.GameLog(ctx, 1)
		require.NoError(t, err)
		require.Len(t, log, 2)
		assert.Equal(t, &models.PlayerHeader{ID: 2, Name: "2"}, log[0].EnteredBy)
		assert.Nil(t, log[1].EnteredBy)

		require.NoError(t, repo.UndoLastRound(ctx, 1, 2, 6, nil))
		require.NoError(t, repo.UndoLastRound(ctx, 1, 1, 1, nil))

		log, err = repo.GameLog(ctx, 1)
		require.NoError(t, err)
		assert.Empty(t, log)

		undone, err = repo.UndoneEntries(ctx, 1)
		require.NoError(t, err)
		require.Len(t, undone, 2)
		assert.Equal(t, models.Draw, undone[0].Outcome, "in the order they were undone")
		assert.Nil(t, undone[0].EnteredBy)
		assert.Equal(t, &models.PlayerHeader{ID: 6, Name: "6"}, undone[0].UndoneBy)
		assert.Equal(t, [4]int{-16, 48, -16, -16}, undone[0].ScoresBefore, "the self-draw's scores")
		assert.Equal(t, [4]int{-16, 48, -16, -16}, undone[0].ScoresAfter, "a draw changes nothing")

		self := undone[1]
		assert.Equal(t, 1, self.Round)
		assert.Equal(t, models.Self, self.Outcome)
		assert.Equal(t, 2, self.Winner)
		assert.Equal(t, 8, self.Value)
		assert.Equal(t, &models.PlayerHeader{ID: 2, Name: "2"}, self.EnteredBy)
		assert.Equal(t, &models.PlayerHeader{ID: 1, Name: "1"}, self.UndoneBy)
		assert.Equal(t, [4]int{}, self.ScoresBefore, "the first entry starts from zero")
		assert.Equal(t, [4]int{-16, 48, -16, -16}, self.ScoresAfter, "by seat: east, south, west, north")
		assert.True(t, self.EnteredAt.After(before))
		assert.False(t, self.UndoneAt.Before(self.EnteredAt))
	})
}
