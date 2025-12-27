package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGameResult(t *testing.T) {
	testData := map[string]struct {
		gameID int

		mockSetup func(*MockGameRepo, int)

		expectedErr    assert.ErrorAssertionFunc
		expectedResult models.GameResult
	}{
		"ok sorted": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("GameResult", mock.Anything, gameID).Return(
					models.GameResult{
						ID: 100,
						PlayerResults: []models.PlayerResult{
							{ID: 2, Name: "Bob", GP: 1, TP: 2},
							{ID: 3, Name: "Charlie", GP: -1, TP: 1},
							{ID: 1, Name: "Alice", GP: 2, TP: 4},
							{ID: 4, Name: "David", GP: -2, TP: 0},
						},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			expectedResult: models.GameResult{
				ID: 100,
				PlayerResults: []models.PlayerResult{
					{ID: 1, Name: "Alice", GP: 2, TP: 4},
					{ID: 2, Name: "Bob", GP: 1, TP: 2},
					{ID: 3, Name: "Charlie", GP: -1, TP: 1},
					{ID: 4, Name: "David", GP: -2, TP: 0},
				},
			},
		},
		"repo error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("GameResult", mock.Anything, gameID).Return(
					models.GameResult{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"empty result": {
			gameID: 200,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("GameResult", mock.Anything, gameID).Return(
					models.GameResult{
						ID:            200,
						PlayerResults: []models.PlayerResult{},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			expectedResult: models.GameResult{
				ID:            200,
				PlayerResults: []models.PlayerResult{},
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, testCase.gameID)

			ctx := context.Background()
			result, err := service.GameResult(ctx, testCase.gameID)

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
