package event

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestRatings(t *testing.T) {
	const eventID = 7

	testData := map[string]struct {
		mockSetup           func(*MockEventRepo)
		wantRankedPlayerIDs []int
		expectedErr         assert.ErrorAssertionFunc
	}{
		"empty": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).
					Return([]models.EventPlayerResultRow{}, nil)
			},
			wantRankedPlayerIDs: []int{},
			expectedErr:         assert.NoError,
		},
		"order by TP then GP": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
					{PlayerID: 1, Name: "A", TP: 2, GP: 100},
					{PlayerID: 2, Name: "B", TP: 4, GP: 50},
					{PlayerID: 3, Name: "C", TP: 2, GP: 120},
				}, nil)
			},
			wantRankedPlayerIDs: []int{2, 3, 1},
			expectedErr:         assert.NoError,
		},
		"tiebreak better session": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
					{PlayerID: 1, Name: "A", TP: 2, GP: 190},
					{PlayerID: 1, Name: "A", TP: 2, GP: 10},
					{PlayerID: 2, Name: "B", TP: 2, GP: 150},
					{PlayerID: 2, Name: "B", TP: 2, GP: 50},
				}, nil)
			},
			wantRankedPlayerIDs: []int{1, 2},
			expectedErr:         assert.NoError,
		},
		"tiebreak equal session GP more games": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
					{PlayerID: 1, Name: "A", TP: 2, GP: 190},
					{PlayerID: 1, Name: "A", TP: 0, GP: 0},
					{PlayerID: 2, Name: "B", TP: 2, GP: 190},
				}, nil)
			},
			wantRankedPlayerIDs: []int{1, 2},
			expectedErr:         assert.NoError,
		},
		"full tiebreak lottery": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
					{PlayerID: 10, Name: "X", TP: 4, GP: 100},
					{PlayerID: 3, Name: "Y", TP: 4, GP: 100},
				}, nil)
			},
			wantRankedPlayerIDs: []int{3, 10},
			expectedErr:         assert.NoError,
		},
		"tiebreaksecond best session": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
					{PlayerID: 1, Name: "A", TP: 2, GP: 100},
					{PlayerID: 1, Name: "A", TP: 2, GP: 50},
					{PlayerID: 1, Name: "A", TP: 2, GP: 50},
					{PlayerID: 2, Name: "B", TP: 2, GP: 100},
					{PlayerID: 2, Name: "B", TP: 2, GP: 80},
					{PlayerID: 2, Name: "B", TP: 2, GP: 20},
				}, nil)
			},
			wantRankedPlayerIDs: []int{2, 1},
			expectedErr:         assert.NoError,
		},
		"repo error": {
			mockSetup: func(m *MockEventRepo) {
				m.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).
					Return(nil, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, eventID).
				Return(models.Event{ID: eventID, Settings: models.DefaultEventSettings()}, nil)
			testCase.mockSetup(repo)

			service := NewService(repo, nil, nil)
			got, err := service.Ratings(context.Background(), eventID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.wantRankedPlayerIDs, rankedPlayerIDs(got.Ratings))
			assert.Empty(t, got.BelowMinimum)
		})
	}
}

func rankedPlayerIDs(entries []models.EventRatingEntry) []int {
	ids := make([]int, len(entries))
	for i, e := range entries {
		ids[i] = e.PlayerID
	}
	return ids
}

// TestRatings_MinGames checks that players below the event's minimum games
// get no place and are listed separately, in the same order.
func TestRatings_MinGames(t *testing.T) {
	const eventID = 7
	settings := models.DefaultEventSettings()
	settings.MinRatingGames = 2

	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{ID: eventID, Settings: settings}, nil)
	repo.EXPECT().ListEndedPlayerResults(mock.Anything, eventID).Return([]models.EventPlayerResultRow{
		{PlayerID: 1, Name: "A", TP: 4, GP: 100}, // one game, best TP
		{PlayerID: 2, Name: "B", TP: 1, GP: 10},
		{PlayerID: 2, Name: "B", TP: 1, GP: 10},
		{PlayerID: 3, Name: "C", TP: 2, GP: 50},
		{PlayerID: 3, Name: "C", TP: 2, GP: 50},
		{PlayerID: 4, Name: "D", TP: 3, GP: 0}, // one game
	}, nil)

	got, err := NewService(repo, nil, nil).Ratings(context.Background(), eventID)

	require.NoError(t, err)
	assert.Equal(t, 2, got.MinGames)
	assert.Equal(t, []models.EventRatingEntry{
		{Rank: 1, PlayerID: 3, Name: "C", Games: 2, TP: 4, GP: 100},
		{Rank: 2, PlayerID: 2, Name: "B", Games: 2, TP: 2, GP: 20},
	}, got.Ratings)
	assert.Equal(t, []models.EventRatingEntry{
		{PlayerID: 1, Name: "A", Games: 1, TP: 4, GP: 100},
		{PlayerID: 4, Name: "D", Games: 1, TP: 3, GP: 0},
	}, got.BelowMinimum)
}
