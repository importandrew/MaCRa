package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFalseMahjong(t *testing.T) {
	testData := map[string]struct {
		gameID         int
		playerID       int
		value          int
		forceNextRound bool
		setupGame      func(*Service, *testing.T) int

		mockSetup func(*MockGameRepo, int)

		expectedErr assert.ErrorAssertionFunc
	}{
		"ok with value 30": {
			gameID:         100,
			playerID:       1,
			value:          30,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().FalseMahjong(1, 30, false).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.FalseMahjong &&
						data.LoserID == 1 &&
						data.Value == 30 &&
						data.ForceNextRound == false
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"ok with value 60": {
			gameID:         100,
			playerID:       2,
			value:          60,
			forceNextRound: true,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().FalseMahjong(2, 60, true).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.FalseMahjong &&
						data.LoserID == 2 &&
						data.Value == 60 &&
						data.ForceNextRound == true
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"game not found": {
			gameID:         999,
			playerID:       1,
			value:          30,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				// Don't add game
				return 999
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// No mock setup needed
			},
			expectedErr: assert.Error,
		},
		"invalid value": {
			gameID:         100,
			playerID:       1,
			value:          50, // Invalid, must be 30 or 60
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().FalseMahjong(1, 50, false).Return(errors.New("invalid value"))
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// No mock setup needed as game engine will return error
			},
			expectedErr: assert.Error,
		},
		"repo error": {
			gameID:         100,
			playerID:       1,
			value:          30,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().FalseMahjong(1, 30, false).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.Anything).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
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
			err := service.FalseMahjong(
				ctx,
				testCase.gameID,
				testCase.playerID,
				testCase.value,
				testCase.forceNextRound,
			)

			testCase.expectedErr(t, err)
			mockGameRepo.AssertExpectations(t)
		})
	}
}
