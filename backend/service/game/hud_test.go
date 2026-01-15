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
		gameID int
		linger int

		mockSetup func(*MockGameRepo, *MockPlayerRepo, int)

		setupGame func(*Service, *testing.T) int

		expectedErr assert.ErrorAssertionFunc
		checkHUD    func(*testing.T, models.GameHUD)
	}{
		"ok final game": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
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
			setupGame: func(s *Service, t *testing.T) int {
				// No live game needed for ended games
				return 100
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
				mg.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
				// Logs calls Players
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				).Times(2) // Called by Logs and Players
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil).Times(2)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil).Times(2)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil).Times(2)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil).Times(2)
				// Logs calls GameLog - empty logs means no diff/message
				mg.On("GameLog", mock.Anything, gameID).Return([]models.RoundLog{}, nil)
			},
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 10},
					2: {GP: -5},
					3: {GP: 0},
					4: {GP: -5},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				// StartWind is only called when showDiffAndMessage is true (when there are round logs)
				// Since we have empty logs, StartWind won't be called
				s.liveGames[100] = mockGame
				return 100
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
				mg.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				).Times(2)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil).Times(2)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil).Times(2)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil).Times(2)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil).Times(2)
				// GameLog with recent round log
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
							EastScore:     32,
							SouthScore:    -16,
							WestScore:     -8,
							NorthScore:    -8,
						},
					},
					nil,
				)
			},
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 32},
					2: {GP: -16},
					3: {GP: -8},
					4: {GP: -8},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().StartWind(1).Return(models.East)
				mockGame.EXPECT().StartWind(2).Return(models.South)
				mockGame.EXPECT().StartWind(3).Return(models.West)
				mockGame.EXPECT().StartWind(4).Return(models.North)
				s.liveGames[100] = mockGame
				return 100
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
			setupGame: func(s *Service, t *testing.T) int {
				return 100
			},
			expectedErr: assert.Error,
		},
		"unexpected status": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameCreated, nil)
			},
			setupGame: func(s *Service, t *testing.T) int {
				return 100
			},
			expectedErr: assert.Error,
		},
		"final game result error": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
				mg.On("Status", mock.Anything, gameID).Return(models.GameEnded, nil)
				mg.On("GameResult", mock.Anything, gameID).Return(
					models.GameResult{},
					errors.New("database error"),
				)
			},
			setupGame: func(s *Service, t *testing.T) int {
				return 100
			},
			expectedErr: assert.Error,
		},
		"live game logs error": {
			gameID: 100,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
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
				mg.On("GameLog", mock.Anything, gameID).Return(nil, errors.New("database error"))
			},
			setupGame: func(s *Service, t *testing.T) int {
				return 100
			},
			expectedErr: assert.Error,
		},
		"live game not found": {
			gameID: 999,
			linger: 0,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Status", mock.Anything, gameID).Return(models.GameLive, nil)
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
				mg.On("GameLog", mock.Anything, gameID).Return([]models.RoundLog{}, nil)
			},
			setupGame: func(s *Service, t *testing.T) int {
				// Don't add any games
				return 999
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
			testCase.mockSetup(mockGameRepo, mockPlayerRepo, gameID)

			ctx := context.Background()
			hud, err := service.HUD(ctx, gameID, testCase.linger)

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
