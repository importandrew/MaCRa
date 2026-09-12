package event

import (
	"context"
	"fmt"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, db database.DB, userID int) {
	t.Helper()

	ctx := context.Background()
	conn, err := db.Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		INSERT INTO users(id, name, email, password)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING`,
		userID, fmt.Sprintf("user%d", userID), fmt.Sprintf("%d@test.com", userID), "password123",
	)
	require.NoError(t, err)
}

func TestCreate_regularEventSeedsOwnerRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "regular", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		ownerID := 42
		createTestUser(t, db, ownerID)

		eventID, err := repo.Create(ctx, models.Event{
			Type:         models.EventRegular,
			Name:         "Regular",
			OwnerUserID:  ownerID,
			Visibility:   models.EventVisibilityPrivate,
			RatingsReach: models.ContentReachMembers,
			GamesReach:   models.ContentReachMembers,
		})
		require.NoError(t, err)

		roles, err := repo.GetMemberRoles(ctx, eventID, ownerID)
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.EventRole{models.EventPlayer, models.EventReferee}, roles)
	})
}

func TestCreate_tournamentEventNoOwnerRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tournament", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		ownerID := 43
		createTestUser(t, db, ownerID)

		eventID, err := repo.Create(ctx, models.Event{
			Type:         models.EventTournament,
			Name:         "Tournament",
			OwnerUserID:  ownerID,
			Visibility:   models.EventVisibilityPrivate,
			RatingsReach: models.ContentReachMembers,
			GamesReach:   models.ContentReachMembers,
		})
		require.NoError(t, err)

		roles, err := repo.GetMemberRoles(ctx, eventID, ownerID)
		require.NoError(t, err)
		assert.Empty(t, roles)
	})
}
