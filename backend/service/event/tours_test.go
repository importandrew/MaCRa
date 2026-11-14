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
	assert.ErrorIs(t, err, models.ErrNoSeatingPlan)
}

func TestImportSeatingPlan(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEventRepo(t)
	service := NewService(repo, nil, nil)

	_, err := service.ImportSeatingPlan(ctx, 1, "1,1,1,2,3\n")
	assert.ErrorIs(t, err, models.ErrSeatingFileRow, "a bad file isn't stored")

	tours := models.PlanTours{3: {{1, 2, 3, 4}}}
	repo.EXPECT().SaveSeatingPlan(ctx, 1, models.PlanImported, tours).Return(nil).Once()
	info, err := service.ImportSeatingPlan(ctx, 1, "3,1,1,2,3,4\n")
	require.NoError(t, err)
	assert.Equal(t, []int{3}, info.Tours)
	assert.Equal(t, models.PlanImported, info.Source)
}

func TestGenerateSeatingPlan(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEventRepo(t)
	input := models.SeatingInput{}
	for n := 1; n <= 4; n++ {
		input.Participants = append(input.Participants, models.Participant{UserID: 100 + n, Number: n})
	}
	repo.EXPECT().SeatingPlanInput(ctx, 1).Return(input, 3, 4, nil).Once()
	tours := models.PlanTours{3: {{1, 2, 3, 4}}, 4: {{1, 2, 3, 4}}}
	repo.EXPECT().SaveSeatingPlan(ctx, 1, models.PlanGenerated, tours).Return(nil).Once()

	info, err := NewService(repo, nil, nil).GenerateSeatingPlan(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []int{3, 4}, info.Tours, "the remaining tours")
	assert.Equal(t, 4, info.Players)
}
