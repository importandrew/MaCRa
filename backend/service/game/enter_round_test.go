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

		mockSetup func(*MockGameRepo, int)

		expectedErr assert.ErrorAssertionFunc
	}{
		"ok deal in": {
			playerID: 1,
			outcome:  "deal-in",
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() to find the player
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				// EnterRound calls EnterRoundResult, Players, Score, Round, ID
				mockGame.EXPECT().EnterRoundResult(8, 2, 1).Return(nil)
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
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Deal &&
						data.WinnerID == 1 &&
						data.LoserID == 2 &&
						data.Value == 8
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"ok self": {
			playerID: 1,
			outcome:  "self",
			value:    8,
			winnerID: 1,
			loserID:  0,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(8, 0, 1).Return(nil)
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
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Self &&
						data.WinnerID == 1 &&
						data.LoserID == 0 &&
						data.Value == 8
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"ok draw": {
			playerID: 1,
			outcome:  "draw",
			value:    0,
			winnerID: 0,
			loserID:  0,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(0, 0, 0).Return(nil)
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
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// Round will be advanced by EnterRoundResult, so expect round + 1
				// Scores will also change, so we just check the essential fields
				m.On("SaveRoundOutcome", mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round > 0 &&
						data.Outcome == models.Draw &&
						data.WinnerID == 0 &&
						data.LoserID == 0 &&
						data.Value == 0
				})).Return(nil)
			},
			expectedErr: assert.NoError,
		},
		"player not in game": {
			playerID: 999,
			outcome:  "deal-in",
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// No mock setup needed
			},
			expectedErr: assert.Error,
		},
		"game engine error invalid hand value": {
			playerID: 1,
			outcome:  "deal-in",
			value:    5, // Less than 8, should fail
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(5, 2, 1).Return(errors.New("hand score should be 8 or higher"))
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(m *MockGameRepo, gameID int) {
				// No mock setup needed as game engine will return error
			},
			expectedErr: assert.Error,
		},
		"repo error": {
			playerID: 1,
			outcome:  "deal-in",
			value:    8,
			winnerID: 1,
			loserID:  2,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				mockGame.EXPECT().EnterRoundResult(8, 2, 1).Return(nil)
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
			err := service.EnterRound(
				ctx,
				testCase.playerID,
				testCase.outcome,
				testCase.value,
				testCase.winnerID,
				testCase.loserID,
			)

			testCase.expectedErr(t, err)
			mockGameRepo.AssertExpectations(t)
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
	mockGame.EXPECT().EnterRoundResult(8, 2, 1).Return(nil)
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

	mockGameRepo.On("SaveRoundOutcome", mock.Anything, 100, mock.MatchedBy(func(data models.RoundData) bool {
		return data.Round == 16 && data.Outcome == models.Deal
	})).Return(nil)

	mockGameRepo.On("Finalize", mock.Anything, mock.MatchedBy(func(data models.FinalGameData) bool {
		return data.GameID == 100 &&
			data.East.ID == 1 &&
			data.South.ID == 2 &&
			data.West.ID == 3 &&
			data.North.ID == 4
	})).Return(nil)

	ctx := context.Background()
	err := service.EnterRound(ctx, 1, "deal-in", 8, 1, 2)

	assert.NoError(t, err)
	mockGameRepo.AssertExpectations(t)
}

func TestEnterRound_FinalRound_Error(t *testing.T) {
	// Test error handling when finalization fails
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	// Create a mock game engine set to round 16
	mockGame := NewMockGameEngine(t)
	mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
	mockGame.EXPECT().EnterRoundResult(8, 2, 1).Return(nil)
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

	mockGameRepo.On("SaveRoundOutcome", mock.Anything, 100, mock.Anything).Return(nil)
	mockGameRepo.On("Finalize", mock.Anything, mock.Anything).Return(errors.New("finalization error"))

	ctx := context.Background()
	err := service.EnterRound(ctx, 1, "deal-in", 8, 1, 2)

	assert.Error(t, err)
	mockGameRepo.AssertExpectations(t)
}
