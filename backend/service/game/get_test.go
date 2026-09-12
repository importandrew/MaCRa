package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLiveGame(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		mockSetup   func(*MockGameRepo, int)
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine)
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLiveGame(m, gameID, nil)
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.NotNil(t, g)
				assert.Equal(t, 100, g.ID())
			},
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

			g, err := service.LiveGame(context.Background(), testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				assert.Nil(t, g)
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g)
			}
		})
	}
}

func TestCurrentGame(t *testing.T) {
	testData := map[string]struct {
		playerID    int
		mockSetup   func(*MockGameRepo)
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine, int)
	}{
		"ok": {
			playerID: 1,
			mockSetup: func(m *MockGameRepo) {
				expectCurrentGame(m, 1, 100, nil, 1)
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine, playerID int) {
				assert.NotNil(t, g)
				players := g.Players()
				found := false
				for _, id := range players {
					if id == playerID {
						found = true
						break
					}
				}
				assert.True(t, found, "player should be in the game")
			},
		},
		"player not in any game": {
			playerID: 999,
			mockSetup: func(m *MockGameRepo) {
				m.EXPECT().CurrentGame(mock.Anything, 999, testEventID).Return(0, models.ErrGameNotFound)
			},
			expectedErr: assert.Error,
		},
		"no live games": {
			playerID: 1,
			mockSetup: func(m *MockGameRepo) {
				m.EXPECT().CurrentGame(mock.Anything, 1, testEventID).Return(0, models.ErrGameNotFound)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo)

			g, err := service.CurrentGame(context.Background(), testCase.playerID, testEventID)

			testCase.expectedErr(t, err)
			if err != nil {
				assert.Nil(t, g)
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g, testCase.playerID)
			}
		})
	}
}
