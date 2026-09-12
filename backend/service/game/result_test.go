package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
				m.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
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
				m.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
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
				m.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
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

func TestGameResults(t *testing.T) {
	t.Run("empty ids", func(t *testing.T) {
		svc := NewService(NewMockGameRepo(t), NewMockPlayerRepo(t))
		got, err := svc.GameResults(context.Background(), nil)
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("sorted per game", func(t *testing.T) {
		mockGameRepo := NewMockGameRepo(t)
		mockGameRepo.EXPECT().GameResults(mock.Anything, []int{100, 200}).Return(
			map[int]models.GameResult{
				100: {
					ID: 100,
					PlayerResults: []models.PlayerResult{
						{ID: 2, Name: "Bob", GP: 1, TP: 2},
						{ID: 1, Name: "Alice", GP: 2, TP: 4},
					},
				},
				200: {
					ID: 200,
					PlayerResults: []models.PlayerResult{
						{ID: 3, Name: "Charlie", GP: 0, TP: 1},
					},
				},
			},
			nil,
		)

		svc := NewService(mockGameRepo, NewMockPlayerRepo(t))
		got, err := svc.GameResults(context.Background(), []int{100, 200})
		assert.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, 1, got[100].PlayerResults[0].ID)
		assert.Equal(t, 2, got[100].PlayerResults[1].ID)
		assert.Equal(t, 3, got[200].PlayerResults[0].ID)
	})

	t.Run("repo error", func(t *testing.T) {
		mockGameRepo := NewMockGameRepo(t)
		mockGameRepo.EXPECT().GameResults(mock.Anything, []int{1}).
			Return(nil, errors.New("database error"))

		svc := NewService(mockGameRepo, NewMockPlayerRepo(t))
		_, err := svc.GameResults(context.Background(), []int{1})
		assert.Error(t, err)
	})
}
