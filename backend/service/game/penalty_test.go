package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPenalty(t *testing.T) {
	testData := map[string]struct {
		gameID         int
		playerID       int
		value          int
		forceNextRound bool
		setupGame      func(*Service, *testing.T) int

		mockSetup func(*MockGameRepo, int)

		expectedErr assert.ErrorAssertionFunc
	}{
		"ok": {
			gameID:         100,
			playerID:       1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Penalty(1, 10, false).Return(nil)
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
					return data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == false
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"ok with force next round": {
			gameID:         100,
			playerID:       2,
			value:          10,
			forceNextRound: true,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Penalty(2, 10, true).Return(nil)
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
					return data.Outcome == models.Penalty &&
						data.LoserID == 2 &&
						data.Value == 10 &&
						data.ForceNextRound == true
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"game not found": {
			gameID:         999,
			playerID:       1,
			value:          10,
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
		"repo error": {
			gameID:         100,
			playerID:       1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Penalty(1, 10, false).Return(nil)
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
			err := service.Penalty(ctx, testCase.gameID, testCase.playerID, testCase.value, testCase.forceNextRound)

			testCase.expectedErr(t, err)
			mockGameRepo.AssertExpectations(t)
		})
	}
}
