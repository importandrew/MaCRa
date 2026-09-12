package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLastGameResult(t *testing.T) {
	testData := map[string]struct {
		playerID int

		mockSetup func(*MockGameRepo, int)

		expectedErr    assert.ErrorAssertionFunc
		expectedResult models.GameResult
	}{
		"ok": {
			playerID: 1,
			mockSetup: func(m *MockGameRepo, playerID int) {
				m.On("LastGame", mock.Anything, playerID, testEventID).Return(100, nil)
				m.On("Status", mock.Anything, 100).Return(models.GameEnded, nil)
				m.On("GameResult", mock.Anything, 100).Return(
					models.GameResult{
						ID: 100,
						PlayerResults: []models.PlayerResult{
							{ID: 1, Name: "Alice", GP: 10, TP: 4},
							{ID: 2, Name: "Bob", GP: -5, TP: 2},
						},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			expectedResult: models.GameResult{
				ID: 100,
				PlayerResults: []models.PlayerResult{
					{ID: 1, Name: "Alice", GP: 10, TP: 4},
					{ID: 2, Name: "Bob", GP: -5, TP: 2},
				},
			},
		},
		"last game error": {
			playerID: 1,
			mockSetup: func(m *MockGameRepo, playerID int) {
				m.On("LastGame", mock.Anything, playerID, testEventID).Return(0, errors.New("game not found"))
			},
			expectedErr: assert.Error,
		},
		"game result error": {
			playerID: 1,
			mockSetup: func(m *MockGameRepo, playerID int) {
				m.On("LastGame", mock.Anything, playerID, testEventID).Return(100, nil)
				m.On("Status", mock.Anything, 100).Return(models.GameEnded, nil)
				m.On("GameResult", mock.Anything, 100).Return(
					models.GameResult{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, testCase.playerID)

			result, err := service.LastGameResult(context.Background(), testCase.playerID, testEventID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedResult.ID, result.ID)
			assert.Len(t, result.PlayerResults, len(testCase.expectedResult.PlayerResults))
			for i, expected := range testCase.expectedResult.PlayerResults {
				assert.Equal(t, expected.ID, result.PlayerResults[i].ID)
				assert.Equal(t, expected.Name, result.PlayerResults[i].Name)
				assert.Equal(t, expected.GP, result.PlayerResults[i].GP)
				assert.Equal(t, expected.TP, result.PlayerResults[i].TP)
			}
		})
	}
}
