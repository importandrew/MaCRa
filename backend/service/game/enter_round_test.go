package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnterRound(t *testing.T) {
	testData := map[string]struct {
		playerID  int
		outcome   models.Outcome
		value     int
		winnerID  int
		loserID   int
		setupGame func(*Service, *testing.T) int

		mockSetup func(*MockGameRepo, *MockPlayerRepo, int)

		expectedErr assert.ErrorAssertionFunc
		checkNow    func(*testing.T, models.GameNow)
	}{
		"ok deal in": {
			playerID: 1,
			outcome:  models.Deal,
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() to find the player
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				// EnterRound calls EnterRoundResult, Players, Score, Round, ID
				mockGame.EXPECT().EnterRoundResult(models.Deal, 8, 2, 1).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 32},
					2: {GP: -16},
					3: {GP: -8},
					4: {GP: -8},
				})
				mockGame.EXPECT().Round().Return(1) // Round advances after EnterRoundResult
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after EnterRound, which calls CurrentGame again
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4}) // CurrentGame in Now()
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 32},
					2: {GP: -16},
					3: {GP: -8},
					4: {GP: -8},
				})
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Deal &&
						data.WinnerID == 1 &&
						data.LoserID == 2 &&
						data.Value == 8
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
		"ok self": {
			playerID: 1,
			outcome:  models.Self,
			value:    8,
			winnerID: 1,
			loserID:  0,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(models.Self, 8, 0, 1).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 48},
					2: {GP: -16},
					3: {GP: -16},
					4: {GP: -16},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after EnterRound
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 48},
					2: {GP: -16},
					3: {GP: -16},
					4: {GP: -16},
				})
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Self &&
						data.WinnerID == 1 &&
						data.LoserID == 0 &&
						data.Value == 8
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
			},
		},
		"ok draw": {
			playerID: 1,
			outcome:  models.Draw,
			value:    0,
			winnerID: 0,
			loserID:  0,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(models.Draw, 0, 0, 0).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
					2: {GP: 0},
					3: {GP: 0},
					4: {GP: 0},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				// Now() is called after EnterRound
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Right(1).Return(2)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 0},
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
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Draw &&
						data.WinnerID == 0 &&
						data.LoserID == 0 &&
						data.Value == 0
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
			},
		},
		"player not in game": {
			playerID: 999,
			outcome:  models.Deal,
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// No mock setup needed
			},
			expectedErr: assert.Error,
		},
		"game engine error invalid hand value": {
			playerID: 1,
			outcome:  models.Deal,
			value:    5, // Less than 8, should fail
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(models.Deal, 5, 2, 1).Return(errors.New("hand score should be 8 or higher"))
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// No mock setup needed as game engine will return error
			},
			expectedErr: assert.Error,
		},
		"repo error": {
			playerID: 1,
			outcome:  models.Deal,
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(models.Deal, 8, 2, 1).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 32},
					2: {GP: -16},
					3: {GP: -8},
					4: {GP: -8},
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
		},
		"round 16": {
			playerID: 1,
			outcome:  models.Deal,
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(models.Deal, 8, 2, 1).Return(nil)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().Status().Return(models.GameEnded)
				score := models.GameScore{
					1: {GP: -2, TP: 0},
					2: {GP: -1, TP: 0},
					3: {GP: 1, TP: 0},
					4: {GP: 2, TP: 0},
				}
				mockGame.EXPECT().Score().Return(score)
				mockGame.EXPECT().Round().Return(16)
				mockGame.EXPECT().ID().Return(100)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round == 16 &&
						data.Outcome == models.Deal &&
						data.WinnerID == 1 &&
						data.LoserID == 2 &&
						data.Value == 8
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
			now, err := service.EnterRound(
				ctx,
				testCase.playerID,
				testCase.outcome,
				testCase.value,
				testCase.winnerID,
				testCase.loserID,
			)

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

func TestEnterRound_FinalRound_OK(t *testing.T) {
	// Test the special case where round 16 triggers finalization
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	// Create a mock game engine set to round 16
	mockGame := NewMockGameEngine(t)
	mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
	mockGame.EXPECT().EnterRoundResult(models.Deal, 8, 2, 1).Return(nil)
	mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
	mockGame.EXPECT().Status().Return(models.GameEnded)
	score := models.GameScore{
		1: {GP: -2, TP: 0},
		2: {GP: -1, TP: 0},
		3: {GP: 1, TP: 0},
		4: {GP: 2, TP: 0},
	}
	mockGame.EXPECT().Score().Return(score)
	mockGame.EXPECT().Round().Return(16)
	mockGame.EXPECT().ID().Return(100)
	service.liveGames[100] = mockGame

	mockGameRepo.EXPECT().SaveRoundOutcome(mock.Anything, 100, mock.MatchedBy(func(data models.RoundData) bool {
		return data.Round == 16 && data.Outcome == models.Deal
	})).Return(nil)

	mockGameRepo.EXPECT().Finalize(mock.Anything, mock.MatchedBy(func(data models.FinalGameData) bool {
		return data.GameID == 100 &&
			data.East.ID == 1 &&
			data.South.ID == 2 &&
			data.West.ID == 3 &&
			data.North.ID == 4
	})).Return(nil)

	ctx := context.Background()
	now, err := service.EnterRound(ctx, 1, models.Deal, 8, 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, 100, now.ID)
	assert.Equal(t, models.GameEnded, now.Status)
}

func TestEnterRound_FinalRound_Error(t *testing.T) {
	// Test error handling when finalization fails
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	// Create a mock game engine set to round 16
	mockGame := NewMockGameEngine(t)
	mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
	mockGame.EXPECT().EnterRoundResult(models.Deal, 8, 2, 1).Return(nil)
	mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
	mockGame.EXPECT().Status().Return(models.GameEnded)
	score := models.GameScore{
		1: {GP: -2, TP: 0},
		2: {GP: -1, TP: 0},
		3: {GP: 1, TP: 0},
		4: {GP: 2, TP: 0},
	}
	mockGame.EXPECT().Score().Return(score)
	mockGame.EXPECT().Round().Return(16)
	mockGame.EXPECT().ID().Return(100)
	service.liveGames[100] = mockGame

	mockGameRepo.EXPECT().SaveRoundOutcome(mock.Anything, 100, mock.Anything).Return(nil)
	mockGameRepo.EXPECT().Finalize(mock.Anything, mock.Anything).Return(errors.New("finalization error"))

	ctx := context.Background()
	_, err := service.EnterRound(ctx, 1, models.Deal, 8, 1, 2)

	assert.Error(t, err)
}
