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

func TestSeatSitter(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "seat sitter", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		games := game.NewRepo(db)
		eventID := createTestTournament(t, db) // players 2–5
		createTestUser(t, db, 6)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventSubstitute}))
		createTestUser(t, db, 7)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 7, []models.EventRole{models.EventObserver}))
		require.NoError(t, repo.Draw(ctx, eventID))

		_, err := generateAndPublish(ctx, repo, eventID)
		require.NoError(t, err)
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		gameID := state.Tables[0].GameID
		seat := state.Tables[0].Seats[1].UserID // south
		change := models.SeatChange{EventID: eventID, GameID: gameID, SeatPlayerID: seat, SitterID: 6, Forfeit: true}

		assert.ErrorIs(t, games.SeatSitter(ctx, models.SeatChange{EventID: eventID, GameID: gameID, SeatPlayerID: seat, SitterID: 7}),
			models.ErrNotSubstitute, "an observer isn't a substitute")
		wrongLog := change
		wrongLog.SeenLogSize = 3
		assert.ErrorIs(t, games.SeatSitter(ctx, wrongLog), models.ErrLogChanged)
		require.NoError(t, games.SeatSitter(ctx, change))
		assert.ErrorIs(t, games.SeatSitter(ctx, change), models.ErrAlreadySeated)
		other := change
		other.SeatPlayerID = state.Tables[0].Seats[2].UserID
		assert.ErrorIs(t, games.SeatSitter(ctx, other), models.ErrPlayerInAnotherGame, "one seat at a time")

		// The substitute sits there: the table is theirs to play; the
		// player's game is still theirs to see.
		seated, err := games.SeatedGame(ctx, 6, eventID)
		require.NoError(t, err)
		assert.Equal(t, gameID, seated)
		_, err = games.SeatedGame(ctx, seat, eventID)
		assert.ErrorIs(t, err, models.ErrGameNotFound)
		sits, err := games.IsPlayerInGame(ctx, gameID, seat)
		require.NoError(t, err)
		assert.False(t, sits)
		place, err := games.HasPlace(ctx, gameID, seat)
		require.NoError(t, err)
		assert.True(t, place)
		at, err := games.SeatOf(ctx, gameID, 6)
		require.NoError(t, err)
		assert.Equal(t, seat, at)

		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		south := state.Tables[0].Seats[1]
		require.NotNil(t, south.Sitter)
		assert.Equal(t, 6, south.Sitter.ID)
		assert.True(t, south.Forfeit)
		assert.Equal(t, []models.Substitute{{ID: 6, Name: south.Sitter.Name, Table: 1}}, state.Substitutes)

		// The player is back: no forfeit any more.
		require.NoError(t, repo.StartTour(ctx, eventID))
		require.NoError(t, games.SeatSitter(ctx, models.SeatChange{EventID: eventID, GameID: gameID, SeatPlayerID: seat, SitterID: seat}))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Nil(t, state.Tables[0].Seats[1].Sitter)
		assert.False(t, state.Tables[0].Seats[1].Forfeit)
		current, err := games.CurrentGame(ctx, seat, eventID)
		require.NoError(t, err)
		assert.Equal(t, gameID, current)

		// The substitute again, forfeiting: the seat's result counts for
		// nobody once the tour is finished.
		require.NoError(t, games.SeatSitter(ctx, change))
		seats := state.Tables[0].Seats
		result := func(i int, gp int, tp float32) models.FinalPlayerData {
			return models.FinalPlayerData{ID: seats[i].UserID, GP: gp, TP: tp}
		}
		require.NoError(t, games.Finalize(ctx, models.FinalGameData{
			GameID: gameID, East: result(0, 30, 4), South: result(1, 10, 2), West: result(2, -10, 1), North: result(3, -30, 0),
		}))
		assert.ErrorIs(t, games.SeatSitter(ctx, models.SeatChange{EventID: eventID, GameID: gameID, SeatPlayerID: seat, SitterID: seat}),
			models.ErrGameNotLive, "not after the game has ended")
		require.NoError(t, repo.FinishTour(ctx, eventID))
		results, _, err := repo.TourResults(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, results, 3)
		for _, r := range results {
			assert.NotEqual(t, seat, r.PlayerID, "the forfeited seat")
		}
	})
}
