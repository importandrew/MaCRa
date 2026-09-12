package game

import (
	"context"
	"errors"
	"testing"
	"time"

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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				dealLog := models.RoundLog{
					ID: 1, Round: 1, EventSequence: 1,
					Outcome: models.Deal, Winner: 1, Loser: 2, Value: 8,
					CreatedAt: time.Now(),
				}
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{{}, {dealLog}})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Deal &&
						data.WinnerID == 1 &&
						data.LoserID == 2 &&
						data.Value == 8
				})).Return(nil)
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, models.GameLive, now.Status)
				assert.Equal(t, 2, now.Round.Num)
			},
		},
		"ok self": {
			playerID: 1,
			outcome:  models.Self,
			value:    8,
			winnerID: 1,
			loserID:  0,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				selfLog := models.RoundLog{
					ID: 1, Round: 1, EventSequence: 1,
					Outcome: models.Self, Winner: 1, Loser: 0, Value: 8,
					CreatedAt: time.Now(),
				}
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{{}, {selfLog}})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Self &&
						data.WinnerID == 1 &&
						data.LoserID == 0 &&
						data.Value == 8
				})).Return(nil)
				expectPlayerGets(mp)
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				drawLog := models.RoundLog{
					ID: 1, Round: 1, EventSequence: 1,
					Outcome: models.Draw, Winner: 0, Loser: 0, Value: 0,
					CreatedAt: time.Now(),
				}
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{{}, {drawLog}})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Draw &&
						data.WinnerID == 0 &&
						data.LoserID == 0 &&
						data.Value == 0
				})).Return(nil)
				expectPlayerGets(mp)
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().CurrentGame(mock.Anything, 999, testEventID).Return(0, models.ErrGameNotFound)
			},
			expectedErr: assert.Error,
		},
		"game engine error invalid hand value": {
			playerID: 1,
			outcome:  models.Deal,
			value:    5,
			winnerID: 1,
			loserID:  2,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				expectCurrentGame(m, 1, gameID, nil, 1)
			},
			expectedErr: assert.Error,
		},
		"repo error": {
			playerID: 1,
			outcome:  models.Deal,
			value:    8,
			winnerID: 1,
			loserID:  2,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				expectCurrentGame(m, 1, gameID, nil, 1)
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				logs := drawRoundLogs(15)
				expectCurrentGame(m, 1, gameID, logs, 1)
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

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, 100)

			now, err := service.EnterRound(
				context.Background(),
				testCase.playerID,
				testCase.outcome,
				testCase.value,
				testCase.winnerID,
				testCase.loserID,
				testEventID,
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
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	logs := drawRoundLogs(15)
	expectCurrentGame(mockGameRepo, 1, 100, logs, 1)

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

	now, err := service.EnterRound(context.Background(), 1, models.Deal, 8, 1, 2, testEventID)

	assert.NoError(t, err)
	assert.Equal(t, 100, now.ID)
	assert.Equal(t, models.GameEnded, now.Status)
}

func TestEnterRound_FinalRound_Error(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	logs := drawRoundLogs(15)
	expectCurrentGame(mockGameRepo, 1, 100, logs, 1)

	mockGameRepo.EXPECT().SaveRoundOutcome(mock.Anything, 100, mock.Anything).Return(nil)
	mockGameRepo.EXPECT().Finalize(mock.Anything, mock.Anything).Return(errors.New("finalization error"))

	_, err := service.EnterRound(context.Background(), 1, models.Deal, 8, 1, 2, testEventID)

	assert.Error(t, err)
}
