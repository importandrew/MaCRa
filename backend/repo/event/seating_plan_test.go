package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestSeatingPlan(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "seating plan", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestTournament(t, db) // players 2–5

		_, err := repo.SeatingPlan(ctx, eventID)
		assert.ErrorIs(t, err, models.ErrNoSeatingPlan)
		assert.ErrorIs(t, repo.DeleteSeatingPlan(ctx, eventID), models.ErrNoSeatingPlan)
		_, _, _, err = repo.SeatingPlanInput(ctx, eventID)
		assert.ErrorIs(t, err, models.ErrNotDrawn, "a plan is by participant numbers")

		require.NoError(t, repo.Draw(ctx, eventID))
		input, first, last, err := repo.SeatingPlanInput(ctx, eventID)
		require.NoError(t, err)
		assert.Len(t, input.Participants, 4)
		assert.Equal(t, 1, first)
		assert.Equal(t, models.DefaultTours, last)

		tours := models.PlanTours{1: {{4, 3, 2, 1}}, 2: {{1, 2, 3, 4}}}
		require.NoError(t, repo.SaveSeatingPlan(ctx, eventID, models.PlanGenerated, tours))
		plan, err := repo.SeatingPlan(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.PlanGenerated, plan.Source)
		assert.Equal(t, tours, plan.Tours)
		assert.False(t, plan.CreatedAt.IsZero())

		// A tour's seating gets the plan and the tour's number.
		var got models.SeatingInput
		_, err = repo.GenerateTour(ctx, eventID, models.Seating{Method: models.SeatingFixed},
			func(input models.SeatingInput) ([][4]int, error) {
				got = input
				return inOrder(input)
			})
		require.NoError(t, err)
		assert.Equal(t, 1, got.Tour)
		assert.Equal(t, tours, got.Plan)
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, state.Current)
		assert.Equal(t, models.SeatingFixed, state.Current.Seating)

		// A published tour is seated: the plan goes on from the next one,
		// counting its meetings.
		require.NoError(t, repo.PublishTour(ctx, eventID))
		input, first, _, err = repo.SeatingPlanInput(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, 2, first)
		assert.Len(t, input.Meetings, 6, "the published table's pairs")

		// Importing replaces the plan.
		imported := models.PlanTours{2: {{2, 1, 4, 3}}}
		require.NoError(t, repo.SaveSeatingPlan(ctx, eventID, models.PlanImported, imported))
		plan, err = repo.SeatingPlan(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.PlanImported, plan.Source)
		assert.Equal(t, imported, plan.Tours)

		require.NoError(t, repo.DeleteSeatingPlan(ctx, eventID))
		_, err = repo.SeatingPlan(ctx, eventID)
		assert.ErrorIs(t, err, models.ErrNoSeatingPlan)

		regularID := createTestEvent(t, db)
		assert.ErrorIs(t, repo.SaveSeatingPlan(ctx, regularID, models.PlanImported, imported), models.ErrNotTournament)
		_, _, _, err = repo.SeatingPlanInput(ctx, regularID)
		assert.ErrorIs(t, err, models.ErrNotTournament)
	})
}
