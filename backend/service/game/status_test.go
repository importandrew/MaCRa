package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStatus(t *testing.T) {
	testData := map[string]struct {
		gameID int

		mockSetup func(*MockGameRepo, int)

		expectedErr    assert.ErrorAssertionFunc
		expectedStatus models.GameStatus
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
			},
			expectedErr:    assert.NoError,
			expectedStatus: models.GameLive,
		},
		"repo error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("Status", mock.Anything, gameID).Return(models.GameStatus(""), errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"game ended": {
			gameID: 200,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
			},
			expectedErr:    assert.NoError,
			expectedStatus: models.GameEnded,
		},
		"game aborted": {
			gameID: 300,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("Status", mock.Anything, gameID).Return(models.GameAborted, nil)
			},
			expectedErr:    assert.NoError,
			expectedStatus: models.GameAborted,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, testCase.gameID)

			ctx := context.Background()
			status, err := service.Status(ctx, testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedStatus, status)
		})
	}
}
