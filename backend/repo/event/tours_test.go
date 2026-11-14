package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

// inOrder seats the participants as given, four to a table.
func inOrder(participants []models.Participant) ([][4]int, error) {
	tables := [][4]int{}
	for i := 0; i+4 <= len(participants); i += 4 {
		tables = append(tables, [4]int{participants[i].UserID, participants[i+1].UserID,
			participants[i+2].UserID, participants[i+3].UserID})
	}
	return tables, nil
}

func TestTours(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tours", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		games := game.NewRepo(db)
		eventID := createTestTournament(t, db) // players 2–5

		_, err := repo.PublishTour(ctx, eventID, models.SeatingRandom, inOrder)
		assert.ErrorIs(t, err, models.ErrNotDrawn, "the draw comes first")

		createTestUser(t, db, 6)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventPlayer}))
		require.NoError(t, repo.Draw(ctx, eventID))
		_, err = repo.PublishTour(ctx, eventID, models.SeatingRandom, inOrder)
		assert.ErrorIs(t, err, models.ErrTablesIncomplete, "5 players: 3 missing")
		var dataErr *models.DataError
		require.ErrorAs(t, err, &dataErr)
		assert.Equal(t, 3, dataErr.Data["missing"])
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, 5, state.Active)
		assert.Equal(t, 3, state.Missing)

		require.NoError(t, repo.RemoveMembers(ctx, eventID, []int{6}))
		require.NoError(t, repo.Draw(ctx, eventID))
		tour, err := repo.PublishTour(ctx, eventID, models.SeatingRandom, inOrder)
		require.NoError(t, err)
		assert.Equal(t, 1, tour.Number)
		assert.Equal(t, models.TourPublished, tour.Status)

		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, state.Current)
		assert.Equal(t, 0, state.Next, "no next tour while one is published")
		require.Len(t, state.Tables, 1)
		table := state.Tables[0]
		assert.Equal(t, 1, table.Table)
		require.Len(t, table.Seats, 4)
		assert.Equal(t, "east", table.Seats[0].Wind)
		assert.Equal(t, 1, table.Seats[0].Number, "seated in number order by inOrder")

		status, err := games.Status(ctx, table.GameID)
		require.NoError(t, err)
		assert.Equal(t, models.GameCreated, status)
		seated, err := games.SeatedGame(ctx, table.Seats[2].UserID, eventID)
		require.NoError(t, err)
		assert.Equal(t, table.GameID, seated, "a published tour's game is the player's table")
		_, err = games.CurrentGame(ctx, table.Seats[2].UserID, eventID)
		assert.ErrorIs(t, err, models.ErrGameNotFound, "but no live game to enter rounds into")
		seat, err := games.TourSeat(ctx, table.GameID)
		require.NoError(t, err)
		assert.Equal(t, &models.TourSeat{Tour: 1, Table: 1}, seat)

		_, err = repo.PublishTour(ctx, eventID, models.SeatingRandom, inOrder)
		assert.ErrorIs(t, err, models.ErrTourInProgress)
		assert.ErrorIs(t, repo.Draw(ctx, eventID), models.ErrTourInProgress, "cancel the seating before a new draw")

		require.NoError(t, repo.CancelTour(ctx, eventID))
		_, err = games.Status(ctx, table.GameID)
		assert.ErrorIs(t, err, models.ErrGameNotFound, "cancelling deletes the games")
		assert.ErrorIs(t, repo.CancelTour(ctx, eventID), models.ErrTourNotFound)
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Nil(t, state.Current)
		assert.NotNil(t, state.Tables, "an empty list, not null, for the client")
		assert.Equal(t, 1, state.Next)

		// With every tour finished, there is no next one.
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()
		for n := 1; n <= models.DefaultTours; n++ {
			_, err = conn.Exec(ctx, `
				INSERT INTO tours (event_id, number, status, seating, published_at)
				VALUES ($1, $2, 'finished', 'random', now())`, eventID, n)
			require.NoError(t, err)
		}
		_, err = repo.PublishTour(ctx, eventID, models.SeatingRandom, inOrder)
		assert.ErrorIs(t, err, models.ErrAllToursPlayed)
	})
}
