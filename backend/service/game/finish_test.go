package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFinish(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		mockSetup   func(*MockGameRepo, int)
		expectedErr assert.ErrorAssertionFunc
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLiveGame(m, gameID, nil)
				m.EXPECT().Finalize(mock.Anything, mock.AnythingOfType("models.FinalGameData")).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"repo finalize error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLiveGame(m, gameID, nil)
				m.EXPECT().Finalize(mock.Anything, mock.AnythingOfType("models.FinalGameData")).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"game not live": {
			gameID: 999,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Status(mock.Anything, gameID).Return(models.GameEnded, nil)
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

			err := service.Finish(context.Background(), testCase.gameID)
			testCase.expectedErr(t, err)
		})
	}
}
