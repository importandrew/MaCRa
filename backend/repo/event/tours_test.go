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

// inOrder seats the participants as given, four to a table.
func inOrder(participants []models.Participant) ([][4]int, error) {
	tables := [][4]int{}
	for i := 0; i+4 <= len(participants); i += 4 {
		tables = append(tables, [4]int{participants[i].UserID, participants[i+1].UserID,
			participants[i+2].UserID, participants[i+3].UserID})
	}
	return tables, nil
}

// generateAndPublish drafts the next tour seated by inOrder and publishes
// it, returning the tour as drafted.
func generateAndPublish(ctx context.Context, repo *Repo, eventID int) (models.Tour, error) {
	tour, err := repo.GenerateTour(ctx, eventID, models.SeatingRandom, inOrder)
	if err != nil {
		return tour, err
	}
	return tour, repo.PublishTour(ctx, eventID)
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

		_, err := generateAndPublish(ctx, repo, eventID)
		assert.ErrorIs(t, err, models.ErrNotDrawn, "the draw comes first")

		createTestUser(t, db, 6)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventPlayer}))
		require.NoError(t, repo.Draw(ctx, eventID))
		_, err = generateAndPublish(ctx, repo, eventID)
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
		// A draft is seen by the staff only, and may be generated again.
		assert.ErrorIs(t, repo.PublishTour(ctx, eventID), models.ErrTourNotFound, "nothing drafted")
		draft, err := repo.GenerateTour(ctx, eventID, models.SeatingRandom, inOrder)
		require.NoError(t, err)
		assert.Equal(t, models.TourDraft, draft.Status)
		assert.Nil(t, draft.PublishedAt)
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, state.Tables, 1)
		draftGame := state.Tables[0].GameID
		_, err = games.SeatedGame(ctx, state.Tables[0].Seats[0].UserID, eventID)
		assert.ErrorIs(t, err, models.ErrGameNotFound, "players don't see a draft")
		assert.ErrorIs(t, repo.Draw(ctx, eventID), models.ErrTourInProgress, "no draw over a draft")
		assert.ErrorIs(t, repo.StartTour(ctx, eventID), models.ErrTourNotFound, "a draft can't start")
		_, err = repo.GenerateTour(ctx, eventID, models.SeatingRandom, inOrder)
		require.NoError(t, err, "generated again")
		_, err = games.Status(ctx, draftGame)
		assert.ErrorIs(t, err, models.ErrGameNotFound, "the old draft's games are gone")

		require.NoError(t, repo.PublishTour(ctx, eventID))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, state.Current)
		assert.Equal(t, 1, state.Current.Number)
		assert.Equal(t, models.TourPublished, state.Current.Status)
		assert.NotNil(t, state.Current.PublishedAt)
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

		_, err = generateAndPublish(ctx, repo, eventID)
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
		_, err = generateAndPublish(ctx, repo, eventID)
		assert.ErrorIs(t, err, models.ErrAllToursPlayed)
	})
}

func TestTourRun(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tour run", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		games := game.NewRepo(db)
		eventID := createTestTournament(t, db) // players 2–5
		require.NoError(t, repo.Draw(ctx, eventID))

		assert.ErrorIs(t, repo.StartTour(ctx, eventID), models.ErrTourNotFound, "nothing published")
		_, err := generateAndPublish(ctx, repo, eventID)
		require.NoError(t, err)
		assert.ErrorIs(t, repo.FinishTour(ctx, eventID), models.ErrTourNotFound, "not started")

		require.NoError(t, repo.StartTour(ctx, eventID))
		event, err := repo.Get(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.EventInProgress, event.Status(), "the first tour starts the tournament")
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.TourLive, state.Current.Status)
		gameID := state.Tables[0].GameID
		assert.Equal(t, models.GameLive, state.Tables[0].Status)
		assert.ErrorIs(t, repo.CancelTour(ctx, eventID), models.ErrTourNotFound, "a started tour can't be cancelled")

		err = repo.FinishTour(ctx, eventID)
		assert.ErrorIs(t, err, models.ErrTablesStillPlaying)
		var dataErr *models.DataError
		require.ErrorAs(t, err, &dataErr)
		assert.Equal(t, 1, dataErr.Data["playing"])

		seats := state.Tables[0].Seats
		result := func(i int, gp int, tp float32) models.FinalPlayerData {
			return models.FinalPlayerData{ID: seats[i].UserID, GP: gp, TP: tp}
		}
		require.NoError(t, games.Finalize(ctx, models.FinalGameData{
			GameID: gameID, East: result(0, 30, 4), South: result(1, 10, 2), West: result(2, -10, 1), North: result(3, -30, 0),
		}))
		rating, err := repo.ListEndedPlayerResults(ctx, eventID)
		require.NoError(t, err)
		assert.Empty(t, rating, "a tour's results count once it is finished")

		require.NoError(t, repo.FinishTour(ctx, eventID))
		rating, err = repo.ListEndedPlayerResults(ctx, eventID)
		require.NoError(t, err)
		assert.Len(t, rating, 4)

		// Tour 2 is aborted: its games don't count, the number is played again.
		_, err = generateAndPublish(ctx, repo, eventID)
		require.NoError(t, err)
		require.NoError(t, repo.StartTour(ctx, eventID))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, 2, state.Current.Number)
		aborted := state.Tables[0].GameID
		require.NoError(t, repo.AbortTour(ctx, eventID))
		status, err := games.Status(ctx, aborted)
		require.NoError(t, err)
		assert.Equal(t, models.GameAborted, status)
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, 1, state.Current.Number, "the current tour is the last one not aborted")
		assert.Equal(t, 2, state.Next, "the aborted number is played again")
		_, err = generateAndPublish(ctx, repo, eventID)
		require.NoError(t, err, "tour 2 again")
	})
}

func TestTourTimer(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tour timer", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		games := game.NewRepo(db)
		eventID := createTestTournament(t, db)
		require.NoError(t, repo.Draw(ctx, eventID))
		_, err := generateAndPublish(ctx, repo, eventID)
		require.NoError(t, err)
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		gameID := state.Tables[0].GameID

		limit, err := games.TimeLimit(ctx, gameID)
		require.NoError(t, err)
		assert.Nil(t, limit, "no time before the start")
		assert.ErrorIs(t, repo.AddTime(ctx, eventID, gameID, 5), models.ErrTourNotFound, "no live tour")

		require.NoError(t, repo.StartTour(ctx, eventID))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, state.Current.Minutes)
		assert.Equal(t, models.DefaultTourMinutes, *state.Current.Minutes, "the tour keeps its length")
		limit, err = games.TimeLimit(ctx, gameID)
		require.NoError(t, err)
		require.NotNil(t, limit)
		start := limit.EndsAt.Add(-models.DefaultTourMinutes * time.Minute)
		assert.Equal(t, models.TimeUpFinishHand, limit.TimeUp)

		require.NoError(t, repo.AddTime(ctx, eventID, gameID, 5))
		limit, err = games.TimeLimit(ctx, gameID)
		require.NoError(t, err)
		assert.WithinDuration(t, start.Add((models.DefaultTourMinutes+5)*time.Minute), limit.EndsAt, time.Second)
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, 5, state.Tables[0].ExtraMinutes)

		// The time-up rule is frozen once the tournament has started; the
		// number of tours can't go below the tours started.
		settings := models.DefaultTournamentSettings()
		settings.TimeUp = models.TimeUpStopNow
		assert.ErrorIs(t, repo.UpdateTournament(ctx, eventID, settings), models.ErrSettingFrozen)
		settings = models.DefaultTournamentSettings()
		settings.TourMinutes = 90
		settings.Tours = 7
		require.NoError(t, repo.UpdateTournament(ctx, eventID, settings), "length and more tours may change")

		// A stopped game: the time is up with "the game ends at once".
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()
		_, err = conn.Exec(ctx, `UPDATE tours SET started_at = now() - interval '3 hours' WHERE event_id = $1`, eventID)
		require.NoError(t, err)
		summary, err := games.Summary(ctx, gameID)
		require.NoError(t, err)
		assert.False(t, summary.StoppedAtTime, "time is up, but the rule finishes the hand")
		_, err = conn.Exec(ctx, `UPDATE events SET tournament = tournament || '{"time_up": "stop_now"}' WHERE id = $1`, eventID)
		require.NoError(t, err)
		summary, err = games.Summary(ctx, gameID)
		require.NoError(t, err)
		assert.True(t, summary.StoppedAtTime)
		timedOut, err := games.TimedOutGames(ctx, time.Now())
		require.NoError(t, err)
		assert.Equal(t, []int{gameID}, timedOut, "the server finishes it")
		timedOut, err = games.TimedOutGames(ctx, time.Now().Add(-3*time.Hour))
		require.NoError(t, err)
		assert.Empty(t, timedOut, "not before its time")
	})
}

func TestCheckFrozen(t *testing.T) {
	old := models.DefaultTournamentSettings() // 5 tours
	more := old
	more.Tours = 6
	fewer := old
	fewer.Tours = 3
	rule := old
	rule.TimeUp, rule.ExtraHands = models.TimeUpExtraHands, 1

	assert.NoError(t, checkFrozen(old, rule, false, 0), "before the start everything may change")
	assert.ErrorIs(t, checkFrozen(old, rule, true, 1), models.ErrSettingFrozen)
	assert.NoError(t, checkFrozen(old, fewer, true, 3), "down to the tours started")
	assert.ErrorIs(t, checkFrozen(old, fewer, true, 4), models.ErrSettingFrozen, "below the tours started")
	assert.ErrorIs(t, checkFrozen(old, more, true, 5), models.ErrSettingFrozen, "the last tour has started")
	assert.NoError(t, checkFrozen(old, old, true, 5), "unchanged is fine")
}
