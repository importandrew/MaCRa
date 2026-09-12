package game

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoad(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		mockSetup   func(*MockGameRepo, int)
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine)
	}{
		"ok empty game": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLoadGame(m, gameID, nil)
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.Equal(t, 100, g.ID())
				assert.Equal(t, 1, g.Round())
			},
		},
		"ok with round logs": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLoadGame(m, gameID, []models.RoundLog{
					{
						ID:            1,
						Round:         1,
						EventSequence: 1,
						Outcome:       models.Deal,
						Winner:        1,
						Loser:         2,
						Value:         8,
						CreatedAt:     time.Now(),
					},
				})
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.Equal(t, 2, g.Round())
			},
		},
		"players error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Players(mock.Anything, gameID).Return(models.StartingPlayers{}, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"game log error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil)
				m.EXPECT().GameLog(mock.Anything, gameID).Return(nil, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, testCase.gameID)

			g, err := service.load(context.Background(), testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g)
			}
		})
	}
}

func TestLoad_WithDifferentOutcomes(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	expectLoadGame(mockGameRepo, 100, []models.RoundLog{
		{
			ID:             1,
			Round:          1,
			EventSequence:  1,
			Outcome:        models.FalseMahjong,
			Loser:          1,
			Value:          30,
			ForceNextRound: false,
			CreatedAt:      time.Now(),
		},
		{
			ID:             2,
			Round:          2,
			EventSequence:  1,
			Outcome:        models.Penalty,
			Loser:          2,
			Value:          10,
			ForceNextRound: true,
			CreatedAt:      time.Now(),
		},
		{
			ID:            3,
			Round:         3,
			EventSequence: 1,
			Outcome:       models.Deal,
			Winner:        3,
			Loser:         4,
			Value:         8,
			CreatedAt:     time.Now(),
		},
	})

	g, err := service.load(context.Background(), 100)

	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, 3, g.Round())
}

func TestLoad_InvalidLogError(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	expectLoadGame(mockGameRepo, 100, []models.RoundLog{
		{
			ID:            1,
			Round:         1,
			EventSequence: 1,
			Outcome:       models.Deal,
			Winner:        1,
			Loser:         2,
			Value:         5,
			CreatedAt:     time.Now(),
		},
	})

	_, err := service.load(context.Background(), 100)

	assert.Error(t, err)
}
