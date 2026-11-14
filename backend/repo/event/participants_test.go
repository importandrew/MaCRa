package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

// createTestTournament creates a tournament owned by user 1 (no roles) with
// users 2–5 as players.
func createTestTournament(t *testing.T, db database.DB) int {
	t.Helper()
	ctx := context.Background()
	for id := 1; id <= 5; id++ {
		createTestUser(t, db, id)
	}
	settings := models.DefaultTournamentSettings()
	eventID, err := NewRepo(db).Create(ctx, models.Event{
		Type: models.EventTournament, Name: "Cup", OwnerUserID: 1,
		Visibility: models.EventVisibilityPrivate, RatingsReach: models.ContentReachMembers,
		GamesReach: models.ContentReachMembers, Rules: models.StandardGameRules(),
		Settings: models.DefaultEventSettings(), Tournament: &settings,
	})
	require.NoError(t, err)
	for id := 2; id <= 5; id++ {
		require.NoError(t, NewRepo(db).AddMemberRoles(ctx, eventID, id, []models.EventRole{models.EventPlayer}))
	}
	return eventID
}

// startTestTournament marks the tournament's first tour as started, as the
// tours will.
func startTestTournament(t *testing.T, db database.DB, eventID int) {
	t.Helper()
	conn, err := db.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `UPDATE events SET started_at = now() WHERE id = $1`, eventID)
	require.NoError(t, err)
}

func TestParticipants(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "participants", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestTournament(t, db)

		list, err := repo.ListParticipants(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, list.Drawn, "not drawn yet")
		assert.Len(t, list.Participants, 4)

		require.NoError(t, repo.Draw(ctx, eventID))
		list, err = repo.ListParticipants(ctx, eventID)
		require.NoError(t, err)
		require.True(t, list.Drawn)
		numbers := map[int]bool{}
		for i, p := range list.Participants {
			require.NotNil(t, p.Number)
			assert.Equal(t, i+1, *p.Number, "listed by number")
			numbers[p.UserID] = true
		}
		assert.Len(t, numbers, 4, "each player has one number")

		// A new player makes the draw out of date until it is redone.
		createTestUser(t, db, 6)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventPlayer}))
		list, err = repo.ListParticipants(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, list.Drawn)
		for _, p := range list.Participants {
			assert.Nil(t, p.Number, "an out-of-date draw shows no numbers")
		}
		require.NoError(t, repo.RemoveMembers(ctx, eventID, []int{6}), "before the start a player can still be removed")
		require.NoError(t, repo.Draw(ctx, eventID))

		assert.ErrorIs(t, repo.Withdraw(ctx, eventID, 2), models.ErrTournamentNotStarted)

		// Once the first tour has started, the players are closed.
		startTestTournament(t, db, eventID)
		assert.ErrorIs(t, repo.Draw(ctx, eventID), models.ErrParticipantsClosed)
		assert.ErrorIs(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventPlayer}), models.ErrParticipantsClosed)
		assert.ErrorIs(t, repo.RemoveMembers(ctx, eventID, []int{3}), models.ErrParticipantsClosed)
		assert.ErrorIs(t, repo.SetMemberRoles(ctx, eventID, 3, []models.EventRole{models.EventReferee}), models.ErrParticipantsClosed)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventSubstitute}), "substitutes still may join")
		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 3, []models.EventRole{models.EventPlayer, models.EventReferee}), "a player may get staff roles")

		require.NoError(t, repo.Withdraw(ctx, eventID, 2))
		assert.ErrorIs(t, repo.Withdraw(ctx, eventID, 2), models.ErrNotParticipant, "already withdrawn")
		assert.ErrorIs(t, repo.Withdraw(ctx, eventID, 6), models.ErrNotParticipant, "a substitute has no number")
		list, err = repo.ListParticipants(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, list.Drawn, "a withdrawal keeps the draw")
		for _, p := range list.Participants {
			assert.Equal(t, p.UserID == 2, p.Withdrawn)
		}

		regularID := createTestEvent(t, db)
		assert.ErrorIs(t, repo.Draw(ctx, regularID), models.ErrNotTournament)
	})
}

func TestDraw_noPlayers(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "no players", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestTournament(t, db)
		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventObserver}))
		require.NoError(t, repo.RemoveMembers(ctx, eventID, []int{3, 4, 5}))

		assert.ErrorIs(t, repo.Draw(ctx, eventID), models.ErrNoPlayers)
	})
}
