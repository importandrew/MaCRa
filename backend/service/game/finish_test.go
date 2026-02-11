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
		gameID    int
		setupGame func(*Service, *testing.T) (int, *MockGameEngine)

		mockSetup func(*MockGameRepo, int, *MockGameEngine)

		expectedErr assert.ErrorAssertionFunc
		verifyGame  func(*testing.T, *Service, int)
	}{
		"ok": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) (int, *MockGameEngine) {
				mockGame := NewMockGameEngine(t)
				s.liveGames[100] = mockGame
				return 100, mockGame
			},
			mockSetup: func(m *MockGameRepo, gameID int, g *MockGameEngine) {
				score := models.GameScore{
					1: {GP: 25000},
					2: {GP: 25000},
					3: {GP: 25000},
					4: {GP: 25000},
				}
				players := [4]int{1, 2, 3, 4}

				g.EXPECT().End().Return(nil)
				g.EXPECT().Score().Return(score)
				g.EXPECT().Players().Return(players)

				m.EXPECT().Finalize(mock.Anything, mock.AnythingOfType("models.FinalGameData")).Return(nil)
			},
			expectedErr: assert.NoError,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// game should be removed from liveGames by finalize
				_, err := s.GetLiveGame(gameID)
				assert.Error(t, err)
				assert.Equal(t, models.ErrGameNotLive, err)
			},
		},
		"end error": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) (int, *MockGameEngine) {
				mockGame := NewMockGameEngine(t)
				s.liveGames[100] = mockGame
				return 100, mockGame
			},
			mockSetup: func(m *MockGameRepo, gameID int, g *MockGameEngine) {
				g.EXPECT().End().Return(errors.New("cannot end game"))
			},
			expectedErr: assert.Error,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// game should still be live if End failed
				_, err := s.GetLiveGame(gameID)
				assert.NoError(t, err)
			},
		},
		"repo finalize error": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) (int, *MockGameEngine) {
				mockGame := NewMockGameEngine(t)
				s.liveGames[100] = mockGame
				return 100, mockGame
			},
			mockSetup: func(m *MockGameRepo, gameID int, g *MockGameEngine) {
				score := models.GameScore{
					1: {GP: 25000},
					2: {GP: 25000},
					3: {GP: 25000},
					4: {GP: 25000},
				}
				players := [4]int{1, 2, 3, 4}

				g.EXPECT().End().Return(nil)
				g.EXPECT().Score().Return(score)
				g.EXPECT().Players().Return(players)

				m.EXPECT().Finalize(mock.Anything, mock.AnythingOfType("models.FinalGameData")).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
				// game should still be live if finalize failed
				_, err := s.GetLiveGame(gameID)
				assert.NoError(t, err)
			},
		},
		"game not in liveGames": {
			gameID: 999,
			setupGame: func(s *Service, t *testing.T) (int, *MockGameEngine) {
				// don't add game
				return 999, nil
			},
			mockSetup:   func(m *MockGameRepo, gameID int, g *MockGameEngine) {},
			expectedErr: assert.Error,
			verifyGame: func(t *testing.T, s *Service, gameID int) {
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

			gameID, liveGame := testCase.setupGame(service, t)
			testCase.mockSetup(mockGameRepo, gameID, liveGame)

			ctx := context.Background()
			err := service.Finish(ctx, gameID)

			testCase.expectedErr(t, err)
			if testCase.verifyGame != nil {
				testCase.verifyGame(t, service, gameID)
			}
		})
	}
}

