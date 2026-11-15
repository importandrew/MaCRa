package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

// planOf returns the plan of tour n in the tours' state.
func planOf(t *testing.T, state models.TourState, n int) models.TourPlan {
	t.Helper()
	require.GreaterOrEqual(t, len(state.Plan), n)
	return state.Plan[n-1]
}

func TestTourPlans(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "tour plans", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestTournament(t, db) // players 2–5
		require.NoError(t, repo.Draw(ctx, eventID))

		// Without a plan every tour is random.
		state, err := repo.TourState(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, state.Plan, models.DefaultTours)
		assert.Equal(t, models.TourPlan{Number: 1, Seating: models.SeatingRandom}, planOf(t, state, 1))

		step := 1
		swiss := models.TourPlan{Number: 2, Seating: models.SeatingSwiss, AvoidRepeats: true}
		require.NoError(t, repo.SetTourPlan(ctx, eventID, swiss))
		require.NoError(t, repo.SetTourPlan(ctx, eventID,
			models.TourPlan{Number: 3, Seating: models.SeatingInterval, Step: &step}))
		err = repo.SetTourPlan(ctx, eventID, models.TourPlan{Number: models.DefaultTours + 1, Seating: models.SeatingRandom})
		assert.ErrorIs(t, err, models.ErrInvalidSeating, "no such tour")
		err = repo.SetTourPlan(ctx, eventID, models.TourPlan{Number: 4, Seating: models.SeatingFixed})
		assert.ErrorIs(t, err, models.ErrFileSeatingMissing, "no tables from a file")

		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, swiss, planOf(t, state, 2))
		assert.Equal(t, &step, planOf(t, state, 3).Step)

		// A file plans its tours fixed and keeps their tables.
		file := models.PlanTours{1: {{4, 3, 2, 1}}, 2: {{1, 2, 3, 4}}}
		err = repo.ImportSeatingFile(ctx, eventID, models.PlanTours{models.DefaultTours + 1: {{1, 2, 3, 4}}})
		assert.ErrorIs(t, err, models.ErrSeatingFileTour)
		require.NoError(t, repo.ImportSeatingFile(ctx, eventID, file))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.TourPlan{Number: 2, Seating: models.SeatingFixed, Tables: file[2], FileTables: 1},
			planOf(t, state, 2), "the Swiss plan gave way to the file")
		assert.Equal(t, models.SeatingInterval, planOf(t, state, 3).Seating, "not in the file")

		// Changing the method keeps the tables, so fixed can come back.
		require.NoError(t, repo.SetTourPlan(ctx, eventID, swiss))
		require.NoError(t, repo.SetTourPlan(ctx, eventID, models.TourPlan{Number: 2, Seating: models.SeatingFixed}))

		// A tour's seating gets its tables and the random tours after it.
		var got models.SeatingInput
		_, err = repo.GenerateTour(ctx, eventID, models.Seating{Method: models.SeatingFixed},
			func(input models.SeatingInput) ([][4]int, error) {
				got = input
				return inOrder(input)
			})
		require.NoError(t, err)
		assert.Equal(t, 1, got.Tour)
		assert.Equal(t, file[1], got.Tables)
		assert.Equal(t, models.DefaultTours-3, got.FutureRandom, "tours 4 on; 2 fixed, 3 interval")
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, planOf(t, state, 1).Seated, "a draft isn't seated yet")

		// Once published, tour 1 is seated: its plan can't change, and a new
		// file leaves its method alone.
		require.NoError(t, repo.PublishTour(ctx, eventID))
		err = repo.SetTourPlan(ctx, eventID, models.TourPlan{Number: 1, Seating: models.SeatingRandom})
		assert.ErrorIs(t, err, models.ErrTourAlreadySeated)
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, planOf(t, state, 1).Seated)

		// Deleting the file: fixed tours not seated become random.
		require.NoError(t, repo.DeleteSeatingFile(ctx, eventID))
		state, err = repo.TourState(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.TourPlan{Number: 2, Seating: models.SeatingRandom}, planOf(t, state, 2))
		assert.Equal(t, models.SeatingFixed, planOf(t, state, 1).Seating, "a seated tour keeps its method")
		assert.Zero(t, planOf(t, state, 1).FileTables)
		assert.True(t, planOf(t, state, 1).Seated)

		regularID := createTestEvent(t, db)
		assert.ErrorIs(t, repo.ImportSeatingFile(ctx, regularID, file), models.ErrNotTournament)
		assert.ErrorIs(t, repo.SetTourPlan(ctx, regularID, swiss), models.ErrNotTournament)
	})
}
