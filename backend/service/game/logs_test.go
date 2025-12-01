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

func TestLogs(t *testing.T) {
	testData := map[string]struct {
		gameID int

		mockSetup func(*MockGameRepo, *MockPlayerRepo, int)

		expectedErr assert.ErrorAssertionFunc
		checkLog    func(*testing.T, models.GameLog)
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				// Setup Players call
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)

				// Setup GameLog call
				mg.On("GameLog", mock.Anything, gameID).Return(
					[]models.RoundLog{
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
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			checkLog: func(t *testing.T, log models.GameLog) {
				assert.Equal(t, 100, log.ID)
				assert.Len(t, log.Players, 4)
				assert.Len(t, log.RoundLogs, 1)
				assert.Equal(t, models.Deal, log.RoundLogs[0].Outcome)
			},
		},
		"players error": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"game log error": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)

				mg.On("GameLog", mock.Anything, gameID).Return(
					nil,
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"empty round logs": {
			gameID: 200,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)

				mg.On("GameLog", mock.Anything, gameID).Return(
					[]models.RoundLog{},
					nil,
				)
			},
			expectedErr: assert.NoError,
			checkLog: func(t *testing.T, log models.GameLog) {
				assert.Equal(t, 200, log.ID)
				assert.Len(t, log.Players, 4)
				assert.Empty(t, log.RoundLogs)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, testCase.gameID)

			ctx := context.Background()
			log, err := service.Logs(ctx, testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				mockGameRepo.AssertExpectations(t)
				mockPlayerRepo.AssertExpectations(t)
				return
			}

			if testCase.checkLog != nil {
				testCase.checkLog(t, log)
			}
			mockGameRepo.AssertExpectations(t)
			mockPlayerRepo.AssertExpectations(t)
		})
	}
}
