package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestUpdateSettings(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "update settings", func(t *testing.T) {
		ctx := context.Background()
		eventID := createTestEvent(t, db)
		settings := models.EventSettings{FinishGame: models.FinishByReferees, AbortGame: models.AbortBySeated, MinRatingGames: 5}

		require.NoError(t, NewRepo(db).UpdateSettings(ctx, eventID, settings))

		event, err := NewRepo(db).Get(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, settings, event.Settings)

		assert.ErrorIs(t, NewRepo(db).UpdateSettings(ctx, eventID+1, settings), models.ErrEventNotFound)
	})
}
