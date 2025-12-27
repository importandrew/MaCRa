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
		setupGame func(*Service, *testing.T) int

		mockSetup func(*MockGameRepo, int)

		expectedErr assert.ErrorAssertionFunc
		verifyGame  func(*testing.T, *Service, int)
	}{
		"ok": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().SetStatus(mock.Anything, gameID, models.GameAborted).Return(nil)
			},
			expectedErr: assert.NoError,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// Verify game is removed from liveGames
				_, err := s.GetLiveGame(gameID)
				assert.Error(t, err)
				assert.Equal(t, models.ErrGameNotLive, err)
			},
		},
		"repo error": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().SetStatus(mock.Anything, gameID, models.GameAborted).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// Game should still be in liveGames if abort failed
				g, err := s.GetLiveGame(gameID)
				assert.NoError(t, err)
				assert.NotNil(t, g)
			},
		},
		"game not in liveGames": {
			gameID: 999,
			setupGame: func(s *Service, t *testing.T) int {
				// Don't add the game
				return 999
			},
			mockSetup:   func(m *MockGameRepo, gameID int) {},
			expectedErr: assert.Error,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// Game should not be in liveGames
				_, err := s.GetLiveGame(gameID)
				assert.Error(t, err)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			gameID := testCase.setupGame(service, t)
			testCase.mockSetup(mockGameRepo, gameID)

			ctx := context.Background()
			err := service.Abort(ctx, gameID)

			testCase.expectedErr(t, err)
			if testCase.verifyGame != nil {
				testCase.verifyGame(t, service, gameID)
			}
		})
	}
}
