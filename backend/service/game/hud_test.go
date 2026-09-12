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

func TestHUD(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		linger      int
		mockSetup   func(*MockGameRepo, *MockPlayerRepo, int)
		expectedErr assert.ErrorAssertionFunc
		checkHUD    func(*testing.T, models.GameHUD)
	}{
		"ok final game": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
				mg.On("GameResult", mock.Anything, gameID).Return(
					models.GameResult{
						ID: 100,
						PlayerResults: []models.PlayerResult{
							{ID: 1, Name: "Alice", GP: 10, TP: 4},
							{ID: 2, Name: "Bob", GP: -5, TP: 2},
						},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			checkHUD: func(t *testing.T, hud models.GameHUD) {
				assert.Nil(t, hud.Live)
				assert.Len(t, hud.Final, 2)
				assert.Equal(t, "Alice", hud.Final[0].Name)
			},
		},
		"ok live game": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil).Times(2)
				mg.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil).Times(3)
				mg.EXPECT().GameLog(mock.Anything, gameID).Return([]models.RoundLog{}, nil).Twice()
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkHUD: func(t *testing.T, hud models.GameHUD) {
				assert.Nil(t, hud.Final)
				assert.NotNil(t, hud.Live)
				assert.Equal(t, 1, hud.Live.Round.Num)
				assert.Len(t, hud.Live.PlayerHUDs, 4)
			},
		},
		"ok live game with message": {
			gameID: 100,
			linger: 30,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				logs := []models.RoundLog{
					{
						ID:            1,
						Round:         1,
						EventSequence: 1,
						Outcome:       models.Deal,
						Winner:        1,
						Loser:         2,
						Value:         8,
						CreatedAt:     time.Now(),
						EastScore:     32,
						SouthScore:    -16,
						WestScore:     -8,
						NorthScore:    -8,
					},
				}
				mg.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil).Times(2)
				mg.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil).Times(3)
				mg.EXPECT().GameLog(mock.Anything, gameID).Return(logs, nil).Twice()
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkHUD: func(t *testing.T, hud models.GameHUD) {
				assert.NotNil(t, hud.Live)
				assert.NotEmpty(t, hud.Live.Message)
				assert.Contains(t, hud.Live.Message, "DEAL-IN")
			},
		},
		"status error": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameStatus(""), errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"unexpected status": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameCreated, nil)
			},
			expectedErr: assert.Error,
		},
		"final game result error": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
				mg.On("GameResult", mock.Anything, gameID).Return(models.GameResult{}, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"live game logs error": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
				mg.On("Players", mock.Anything, gameID).Return(defaultStartingPlayers, nil)
				expectPlayerGets(mp)
				mg.On("GameLog", mock.Anything, gameID).Return(nil, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"live game not found": {
			gameID: 999,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil).Once()
				mg.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil).Twice()
				mg.EXPECT().GameLog(mock.Anything, gameID).Return([]models.RoundLog{}, nil).Once()
				mg.EXPECT().Status(mock.Anything, gameID).Return(models.GameEnded, nil).Once()
				expectPlayerGets(mp)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, testCase.gameID)

			hud, err := service.HUD(context.Background(), testCase.gameID, testCase.linger)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkHUD != nil {
				testCase.checkHUD(t, hud)
			}
		})
	}
}
