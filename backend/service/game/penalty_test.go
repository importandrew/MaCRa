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

func TestPenalty(t *testing.T) {
	testData := map[string]struct {
		playerID       int
		penalizedID    int
		value          int
		forceNextRound bool
		mockSetup      func(*MockGameRepo, *MockPlayerRepo, int)
		expectedErr    assert.ErrorAssertionFunc
		checkNow       func(*testing.T, models.GameNow)
	}{
		"ok": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				penaltyLog := models.RoundLog{
					ID: 1, Round: 1, EventSequence: 1,
					Outcome: models.Penalty, Loser: 1, Value: 10,
					CreatedAt: time.Now(),
				}
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{{}, {penaltyLog}})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == false
				})).Return(nil)
				expectPlayerGets(mp)
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				penaltyLog := models.RoundLog{
					ID: 1, Round: 1, EventSequence: 1,
					Outcome: models.Penalty, Loser: 2, Value: 10, ForceNextRound: true,
					CreatedAt: time.Now(),
				}
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{{}, {penaltyLog}})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Outcome == models.Penalty &&
						data.LoserID == 2 &&
						data.Value == 10 &&
						data.ForceNextRound == true
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
		"game not found": {
			playerID:       999,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				m.EXPECT().CurrentGame(mock.Anything, 999, testEventID).Return(0, models.ErrGameNotFound)
			},
			expectedErr: assert.Error,
		},
		"repo error": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				expectCurrentGame(m, 1, gameID, nil, 1)
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.Anything).Return(errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"round 16 with forceNextRound false": {
			playerID:       1,
			penalizedID:    1,
			value:          10,
			forceNextRound: false,
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				logs := drawRoundLogs(15)
				penaltyLog := models.RoundLog{
					ID: 16, Round: 16, EventSequence: 1,
					Outcome: models.Penalty, Loser: 1, Value: 10,
					CreatedAt: time.Now(),
				}
				secondLoad := append(append([]models.RoundLog{}, logs...), penaltyLog)
				expectCurrentGameLoads(m, 1, gameID, [][]models.RoundLog{logs, secondLoad})
				m.EXPECT().SaveRoundOutcome(mock.Anything, gameID, mock.MatchedBy(func(data models.RoundData) bool {
					return data.Round == 16 &&
						data.Outcome == models.Penalty &&
						data.LoserID == 1 &&
						data.Value == 10 &&
						data.ForceNextRound == false
				})).Return(nil)
				expectPlayerGets(mp)
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
			mockSetup: func(m *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				logs := drawRoundLogs(15)
				expectCurrentGame(m, 1, gameID, logs, 1)
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

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, 100)

			now, err := service.Penalty(
				context.Background(),
				models.Penalty,
				testCase.playerID,
				testCase.penalizedID,
				testCase.value,
				testCase.forceNextRound,
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
