package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestCountPlayers(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "players only", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		for id := 1; id <= 4; id++ {
			createTestUser(t, db, id)
		}

		// A regular event's owner (1) gets the player and referee roles.
		eventID, err := repo.Create(ctx, models.Event{
			Type:         models.EventRegular,
			Name:         "Count",
			OwnerUserID:  1,
			Visibility:   models.EventVisibilityPublic,
			RatingsReach: models.ContentReachPublic,
			GamesReach:   models.ContentReachPublic,
		})
		require.NoError(t, err)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer}))
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 3, []models.EventRole{models.EventPlayer, models.EventObserver}))
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 4, []models.EventRole{models.EventObserver}))

		count, err := repo.CountPlayers(ctx, eventID)

		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})
}
