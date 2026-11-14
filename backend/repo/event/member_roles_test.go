package event

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

func TestSetMemberRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "set roles", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db) // owner 1: player and referee
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer, models.EventReferee}))

		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventReferee, models.EventObserver}))
		roles, err := repo.GetMemberRoles(ctx, eventID, 2)
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.EventRole{models.EventReferee, models.EventObserver}, roles)

		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 1, []models.EventRole{}))
		roles, err = repo.GetMemberRoles(ctx, eventID, 1)
		require.NoError(t, err)
		assert.Empty(t, roles, "the owner may be left without roles")

		roles, err = repo.GetMemberRoles(ctx, eventID, 2)
		require.NoError(t, err)
		assert.Len(t, roles, 2, "other members are untouched")
	})
}

// TestListMembers_ownerWithoutRoles checks that the owner stays in the member
// list after dropping every role.
func TestListMembers_ownerWithoutRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "owner without roles", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db) // owner 1
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer}))
		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 1, []models.EventRole{}))

		members, err := repo.ListMembers(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, []models.EventMember{
			{UserID: 1, Name: "user1", Roles: []models.EventRole{}},
			{UserID: 2, Name: "user2", Roles: []models.EventRole{models.EventPlayer}},
			{UserID: 3, Name: "user3", Roles: []models.EventRole{models.EventPlayer}},
			{UserID: 4, Name: "user4", Roles: []models.EventRole{models.EventPlayer}},
		}, members)
	})
}

// TestMemberChanges_LiveGame checks the refusals for a member in a live game:
// dropping the player role, removing them; and that the owner can't be
// removed.
func TestMemberChanges_LiveGame(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "live game", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db)
		_, err := game.NewRepo(db).NewGame(ctx, [4]int{1, 2, 3, 4}, eventID)
		require.NoError(t, err)

		err = repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventReferee})
		assert.ErrorIs(t, err, models.ErrPlayerRoleInGame)
		assertErrorName(t, err, "user2")
		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer, models.EventReferee}),
			"keeping the player role is fine")

		err = repo.RemoveMembers(ctx, eventID, []int{3})
		assert.ErrorIs(t, err, models.ErrMemberInCurrentGame)
		assert.ErrorIs(t, repo.RemoveMembers(ctx, eventID, []int{1}), models.ErrOwnerNotRemovable)
		assert.ErrorIs(t, repo.SetMemberRoles(ctx, eventID, 99, []models.EventRole{models.EventPlayer}), models.ErrNotMember)
		assert.ErrorIs(t, repo.SetMemberRoles(ctx, eventID, 3, []models.EventRole{}), models.ErrInvalidRole)
	})
}

// TestSetMemberRoles_ConcurrentNewGame checks that dropping a player's role
// and starting a game with them at the same time never leave a live game
// with a player who isn't one: exactly one of the two succeeds.
func TestSetMemberRoles_ConcurrentNewGame(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "concurrent role drop and new game", func(t *testing.T) {
		ctx := context.Background()
		eventID := createTestEvent(t, db)
		events := NewRepo(db)
		games := game.NewRepo(db)

		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for range 20 {
			_, err := conn.Exec(ctx, `DELETE FROM games WHERE event_id = $1`, eventID)
			require.NoError(t, err)
			require.NoError(t, events.AddMemberRoles(ctx, eventID, 4, []models.EventRole{models.EventPlayer}))

			var dropErr, newGameErr error
			var wg sync.WaitGroup
			wg.Go(func() { dropErr = events.SetMemberRoles(ctx, eventID, 4, []models.EventRole{models.EventObserver}) })
			wg.Go(func() { _, newGameErr = games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID) })
			wg.Wait()

			if dropErr == nil {
				assert.ErrorIs(t, newGameErr, models.ErrMissingPlayerRole)
			} else {
				assert.ErrorIs(t, dropErr, models.ErrPlayerRoleInGame)
				assert.NoError(t, newGameErr)
			}
		}
	})
}

// assertErrorName checks the player name sent with err.
func assertErrorName(t *testing.T, err error, name string) {
	t.Helper()
	dataErr, ok := errors.AsType[*models.DataError](err)
	require.True(t, ok, "%v has no data", err)
	assert.Equal(t, name, dataErr.Data["name"])
}

// TestMemberChanges_FinishedEvent checks that a finished event's members
// can't change.
func TestMemberChanges_FinishedEvent(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "finished event", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db)
		require.NoError(t, repo.Finish(ctx, eventID))

		assert.ErrorIs(t, repo.AddMemberRoles(ctx, eventID, 5, []models.EventRole{models.EventPlayer}), models.ErrEventAlreadyFinished)
		assert.ErrorIs(t, repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventReferee}), models.ErrEventAlreadyFinished)
		assert.ErrorIs(t, repo.RemoveMembers(ctx, eventID, []int{2}), models.ErrEventAlreadyFinished)
	})
}
