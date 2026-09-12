package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAbort(t *testing.T) {
	testData := map[string]struct {
		gameID    int
		mockSetup func(*MockGameRepo, int)
		expectedErr assert.ErrorAssertionFunc
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil)
				m.EXPECT().SetStatus(mock.Anything, gameID, models.GameAborted).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"repo error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil)
				m.EXPECT().SetStatus(mock.Anything, gameID, models.GameAborted).Return(errors.New("database error"))
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

			err := service.Abort(context.Background(), testCase.gameID)
			testCase.expectedErr(t, err)
		})
	}
}
