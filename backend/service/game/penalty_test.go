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
		playerID       int
		penalizedID    int
		value          int
		forceNextRound bool
		setupGame      func(*Service, *testing.T) int

		mockSetup func(*MockGameRepo, *MockPlayerRepo, int)

		expectedErr assert.ErrorAssertionFunc
		checkNow    func(*testing.T, models.GameNow)
	}{
		"ok": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Penalty(1, 10, false).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: -10},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after Penalty
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4}) // CurrentGame in Now()
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: -10},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == false
				})).Return(nil)
				// Now() calls playerRepo.Get for each player
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, models.GameLive, now.Status)
				assert.Equal(t, 1, now.Round.Num)
			},
		},
		"ok with force next round": {
			playerID:       1,
			penalizedID:    2,
			value:          10,
			forceNextRound: true,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Penalty(2, 10, true).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: -10},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(2)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after Penalty
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4}) // CurrentGame in Now()
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: -10},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				mockGame.EXPECT().Round().Return(2)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Penalty &&
						data.LoserID == 2 &&
						data.Value == 10 &&
						data.ForceNextRound == true
				})).Return(nil)
				// Now() calls playerRepo.Get for each player
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, models.GameLive, now.Status)
				assert.Equal(t, 2, now.Round.Num)
			},
		},
		"game not found": {
			playerID:       999,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				// Don't add game
				return 999
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// No mock setup needed
			},
			expectedErr: assert.Error,
			checkNow:    nil,
		},
		"repo error": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.Anything).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
			checkNow:    nil,
		},
		"round 16 with forceNextRound false": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame is called first
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Penalty(1, 10, false).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				score := models.GameScore{
					1: {GP: -10, TP: 0},
					2: {GP: 0, TP: 0},
					3: {GP: 0, TP: 0},
					4: {GP: 0, TP: 0},
				}
				mockGame.EXPECT().Score().Return(score)
				mockGame.EXPECT().Round().Return(16)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after Penalty
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4}) // CurrentGame in Now()
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: -10},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				mockGame.EXPECT().Round().Return(16)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round == 16 &&
						data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == false
				})).Return(nil)
				// Now() calls playerRepo.Get for each player
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, models.GameLive, now.Status)
				assert.Equal(t, 16, now.Round.Num)
			},
		},
		"round 16 with forceNextRound true": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: true,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Penalty(1, 10, true).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				score := models.GameScore{
					1: {GP: -10, TP: 0},
					2: {GP: 0, TP: 0},
					3: {GP: 0, TP: 0},
					4: {GP: 0, TP: 0},
				}
				mockGame.EXPECT().Score().Return(score)
				mockGame.EXPECT().Round().Return(16)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameEnded)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round == 16 &&
						data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == true
				})).Return(nil)
				m.EXPECT().Finalize(mock.Anything, mock.MatchedBy(func(data models.FinalGameData) bool {
					return data.GameID == 100 &&
						data.East.ID == 1 &&
						data.South.ID == 2 &&
						data.West.ID == 3 &&
						data.North.ID == 4
				})).Return(nil)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, models.GameEnded, now.Status)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			gameID := testCase.setupGame(service, t)
			testCase.mockSetup(mockGameRepo, mockPlayerRepo, gameID)

			ctx := context.Background()
			now, err := service.Penalty(ctx, models.Penalty, testCase.playerID, testCase.penalizedID, testCase.value, testCase.forceNextRound)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkNow != nil {
				testCase.checkNow(t, now)
			}
		})
	}
}
