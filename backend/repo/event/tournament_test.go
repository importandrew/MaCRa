package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestTournament(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tournament", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		createTestUser(t, db, 1)
		for id := 2; id <= 4; id++ {
			createTestUser(t, db, id)
		}

		defaults := models.DefaultTournamentSettings()
		tournamentID, err := repo.Create(ctx, models.Event{
			Type: models.EventTournament, Name: "Cup", OwnerUserID: 1,
			Visibility: models.EventVisibilityPublic, RatingsReach: models.ContentReachPublic,
			GamesReach: models.ContentReachPublic, Rules: models.StandardGameRules(),
			Settings: models.DefaultEventSettings(), Tournament: &defaults,
		})
		require.NoError(t, err)

		event, err := repo.Get(ctx, tournamentID)
		require.NoError(t, err)
		assert.Equal(t, &defaults, event.Tournament)
		assert.Equal(t, models.EventPreparation, event.Status())
		roles, err := repo.GetMemberRoles(ctx, tournamentID, 1)
		require.NoError(t, err)
		assert.Empty(t, roles, "the owner of a tournament gets no roles")

		changed := models.TournamentSettings{Tours: 7, TourMinutes: 90, TimeUp: models.TimeUpExtraHands, ExtraHands: 1}
		require.NoError(t, repo.UpdateTournament(ctx, tournamentID, changed))
		event, err = repo.Get(ctx, tournamentID)
		require.NoError(t, err)
		assert.Equal(t, &changed, event.Tournament)

		regularID := createTestEvent(t, db)
		assert.ErrorIs(t, repo.UpdateTournament(ctx, regularID, changed), models.ErrNotTournament)
		regular, err := repo.Get(ctx, regularID)
		require.NoError(t, err)
		assert.Nil(t, regular.Tournament)

		// Player and substitute exclude each other, however the roles come.
		require.NoError(t, repo.AddMemberRoles(ctx, tournamentID, 2, []models.EventRole{models.EventSubstitute, models.EventReferee}))
		assert.ErrorIs(t, repo.AddMemberRoles(ctx, tournamentID, 2, []models.EventRole{models.EventPlayer}), models.ErrPlayerAndSubstitute)
		assert.ErrorIs(t, repo.SetMemberRoles(ctx, tournamentID, 2, []models.EventRole{models.EventPlayer, models.EventSubstitute}), models.ErrPlayerAndSubstitute)
		token := "cup-invite"
		require.NoError(t, repo.SetInviteToken(ctx, tournamentID, &token))
		_, err = repo.JoinByInvite(ctx, token, 2)
		assert.ErrorIs(t, err, models.ErrPlayerAndSubstitute, "a substitute can't join as a player by the link")
		roles, err = repo.GetMemberRoles(ctx, tournamentID, 2)
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.EventRole{models.EventSubstitute, models.EventReferee}, roles, "nothing changed")
		require.NoError(t, repo.SetMemberRoles(ctx, tournamentID, 2, []models.EventRole{models.EventPlayer}), "switching substitute to player")

		require.NoError(t, repo.Finish(ctx, tournamentID))
		assert.ErrorIs(t, repo.UpdateTournament(ctx, tournamentID, defaults), models.ErrEventAlreadyFinished)
	})
}
