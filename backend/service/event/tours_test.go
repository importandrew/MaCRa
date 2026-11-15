package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestSeatingFunc_unknown(t *testing.T) {
	_, err := seatingFunc(models.Seating{Method: "swiss-ish"})
	assert.ErrorIs(t, err, models.ErrInvalidSeating)
}

func TestSeatingFunc_intervalStep(t *testing.T) {
	input := models.SeatingInput{}
	for n := 1; n <= 8; n++ { // two tables
		input.Participants = append(input.Participants, models.Participant{UserID: 100 + n, Number: n})
	}
	step := func(k int) *int { return &k }

	_, err := seatingFunc(models.Seating{Method: models.SeatingRandom, Step: step(1)})
	assert.ErrorIs(t, err, models.ErrInvalidSeating, "a step for another method")

	for k, ok := range map[int]bool{0: false, 1: true, 2: true, 3: false} {
		makeTables, err := seatingFunc(models.Seating{Method: models.SeatingInterval, Step: step(k)})
		assert.NoError(t, err)
		_, err = makeTables(input)
		if ok {
			assert.NoError(t, err, "step %d", k)
		} else {
			assert.ErrorIs(t, err, models.ErrInvalidSeating, "step %d", k)
		}
	}
}

func TestSeatingFunc_fixed(t *testing.T) {
	makeTables, err := seatingFunc(models.Seating{Method: models.SeatingFixed})
	require.NoError(t, err)
	_, err = makeTables(models.SeatingInput{Tour: 1})
	assert.ErrorIs(t, err, models.ErrFileSeatingMissing)
}

func TestCheckSeating(t *testing.T) {
	step := 2
	for how, ok := range map[models.Seating]bool{
		{Method: models.SeatingSwiss, AvoidRepeats: true}:                 true,
		{Method: models.SeatingInterval, Step: &step, AvoidRepeats: true}: true,
		{Method: models.SeatingRandom}:                                    true,
		{Method: models.SeatingFixed}:                                     true,
		{Method: models.SeatingRandom, AvoidRepeats: true}:                false,
		{Method: models.SeatingFixed, AvoidRepeats: true}:                 false,
		{Method: models.SeatingSwiss, Step: &step}:                        false,
		{Method: "seat-by-height"}:                                        false,
	} {
		if ok {
			assert.NoError(t, checkSeating(how), "%+v", how)
		} else {
			assert.ErrorIs(t, checkSeating(how), models.ErrInvalidSeating, "%+v", how)
		}
	}
}

func TestSetTourPlan(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEventRepo(t)
	service := NewService(repo, nil, nil)

	zero := 0
	err := service.SetTourPlan(ctx, 1, models.TourPlan{Number: 2, Seating: models.SeatingInterval, Step: &zero})
	assert.ErrorIs(t, err, models.ErrInvalidSeating)
	err = service.SetTourPlan(ctx, 1, models.TourPlan{Number: 2, Seating: models.SeatingRandom, AvoidRepeats: true})
	assert.ErrorIs(t, err, models.ErrInvalidSeating)

	plan := models.TourPlan{Number: 2, Seating: models.SeatingSwiss, AvoidRepeats: true}
	repo.EXPECT().SetTourPlan(ctx, 1, plan).Return(nil).Once()
	require.NoError(t, service.SetTourPlan(ctx, 1, plan))
}

func TestImportSeatingFile(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEventRepo(t)
	service := NewService(repo, nil, nil)

	_, err := service.ImportSeatingFile(ctx, 1, "1,1,1,2,3\n")
	assert.ErrorIs(t, err, models.ErrSeatingFileRow, "a bad file isn't stored")

	tours := models.PlanTours{3: {{1, 2, 3, 4}}}
	repo.EXPECT().ImportSeatingFile(ctx, 1, tours).Return(nil).Once()
	info, err := service.ImportSeatingFile(ctx, 1, "3,1,1,2,3,4\n")
	require.NoError(t, err)
	assert.Equal(t, []int{3}, info.Tours)
}
