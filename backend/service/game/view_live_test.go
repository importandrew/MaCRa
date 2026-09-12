package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestViewLive(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		mockSetup   func(*MockGameRepo, *MockPlayerRepo, int)
		expectedErr assert.ErrorAssertionFunc
		checkView   func(*testing.T, models.GameView)
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				expectLiveGame(mg, gameID, nil)
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkView: func(t *testing.T, view models.GameView) {
				assert.Equal(t, 100, view.ID)
				assert.Equal(t, models.GameLive, view.Status)
				assert.Equal(t, 1, view.Round.Num)
				assert.Len(t, view.PlayerScores, 4)
				assert.Equal(t, 1, view.PlayerScores[0].ID)
				assert.Equal(t, "Alice", view.PlayerScores[0].Name)
				assert.Equal(t, models.East.Kanji(), view.PlayerScores[0].Wind)
				assert.Equal(t, 0, view.PlayerScores[0].Score)
			},
		},
		"game not found": {
			gameID: 999,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.EXPECT().Status(mock.Anything, gameID).Return(models.GameEnded, nil)
			},
			expectedErr: assert.Error,
		},
		"player repo error": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				expectLiveGame(mg, gameID, nil)
				mp.On("Get", mock.Anything, 1).Return(models.Player{}, errors.New("player not found"))
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, testCase.gameID)

			view, err := service.ViewLive(context.Background(), testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkView != nil {
				testCase.checkView(t, view)
			}
		})
	}
}
