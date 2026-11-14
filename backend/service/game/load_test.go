package game

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

func TestLoad(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		mockSetup   func(*MockGameRepo, int)
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine)
	}{
		"ok empty game": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLoadGame(m, gameID, nil)
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.Equal(t, 100, g.ID())
				assert.Equal(t, 1, g.Round())
			},
		},
		"ok with round logs": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				expectLoadGame(m, gameID, []models.RoundLog{
					{
						ID:            1,
						Round:         1,
						EventSequence: 1,
						Outcome:       models.Deal,
						Winner:        1,
						Loser:         2,
						Value:         8,
						CreatedAt:     time.Now(),
					},
				})
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.Equal(t, 2, g.Round())
			},
		},
		"players error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Players(mock.Anything, gameID).Return(models.StartingPlayers{}, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"rules error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil)
				m.EXPECT().Rules(mock.Anything, gameID).Return(models.GameRules{}, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"game log error": {
			gameID: 100,
			mockSetup: func(m *MockGameRepo, gameID int) {
				m.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil)
				m.EXPECT().Rules(mock.Anything, gameID).Return(models.StandardGameRules(), nil)
				m.EXPECT().GameLog(mock.Anything, gameID).Return(nil, errors.New("database error"))
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

			g, _, err := service.load(context.Background(), testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g)
			}
		})
	}
}

func TestLoad_WithDifferentOutcomes(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	expectLoadGame(mockGameRepo, 100, []models.RoundLog{
		{
			ID:             1,
			Round:          1,
			EventSequence:  1,
			Outcome:        models.FalseMahjong,
			Loser:          1,
			Value:          30,
			ForceNextRound: false,
			CreatedAt:      time.Now(),
		},
		{
			ID:             2,
			Round:          2,
			EventSequence:  1,
			Outcome:        models.Penalty,
			Loser:          2,
			Value:          10,
			ForceNextRound: true,
			CreatedAt:      time.Now(),
		},
		{
			ID:            3,
			Round:         3,
			EventSequence: 1,
			Outcome:       models.Deal,
			Winner:        3,
			Loser:         4,
			Value:         8,
			CreatedAt:     time.Now(),
		},
	})

	g, _, err := service.load(context.Background(), 100)

	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, 3, g.Round())
}

func TestLoad_InvalidLogError(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	expectLoadGame(mockGameRepo, 100, []models.RoundLog{
		{
			ID:            1,
			Round:         1,
			EventSequence: 1,
			Outcome:       models.Deal,
			Winner:        1,
			Loser:         2,
			Value:         5,
			CreatedAt:     time.Now(),
		},
	})

	_, _, err := service.load(context.Background(), 100)

	assert.Error(t, err)
}

// TestLoad_ByGameRules checks that a game is replayed by the rules it
// started with: an 8-hand game ends after 8 hands, and its log is valid only
// by its own minimum hand.
func TestLoad_ByGameRules(t *testing.T) {
	short := models.GameRules{Hands: 8, MinHand: 24}

	testData := map[string]struct {
		logs        []models.RoundLog
		expectedErr assert.ErrorAssertionFunc
		wantStatus  models.GameStatus
	}{
		"ends after 8 hands": {
			logs:        drawRoundLogs(8),
			expectedErr: assert.NoError,
			wantStatus:  models.GameEnded,
		},
		"7 hands: still live": {
			logs:        drawRoundLogs(7),
			expectedErr: assert.NoError,
			wantStatus:  models.GameLive,
		},
		"hand below the game's minimum": {
			logs:        []models.RoundLog{{ID: 1, Round: 1, Outcome: models.Self, Winner: 1, Value: 16}},
			expectedErr: assert.Error,
		},
		"hand at the game's minimum": {
			logs:        []models.RoundLog{{ID: 1, Round: 1, Outcome: models.Self, Winner: 1, Value: 24}},
			expectedErr: assert.NoError,
			wantStatus:  models.GameLive,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockGameRepo.EXPECT().Players(mock.Anything, 100).Return(defaultStartingPlayers, nil)
			mockGameRepo.EXPECT().Rules(mock.Anything, 100).Return(short, nil)
			mockGameRepo.EXPECT().GameLog(mock.Anything, 100).Return(testCase.logs, nil)

			g, _, err := NewService(mockGameRepo, NewMockPlayerRepo(t)).load(context.Background(), 100)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.wantStatus, g.Status())
			assert.Equal(t, short, g.Rules())
		})
	}
}
