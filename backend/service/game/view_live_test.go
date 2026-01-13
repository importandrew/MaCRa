package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestViewLive(t *testing.T) {
	testData := map[string]struct {
		gameID int

		setupGame func(*Service, *testing.T) int

		mockSetup func(*MockPlayerRepo, int)

		expectedErr assert.ErrorAssertionFunc
		checkView   func(*testing.T, models.GameView)
	}{
		"ok": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 10},
					2: {GP: -5},
					3: {GP: 0},
					4: {GP: -5},
				})
				mockGame.EXPECT().Round().Return(1)
				// PlayerBySeatInRound is called for each wind in order: East, South, West, North
				mockGame.EXPECT().PlayerBySeatInRound(1, models.East).Return(1)
				mockGame.EXPECT().PlayerBySeatInRound(1, models.South).Return(2)
				mockGame.EXPECT().PlayerBySeatInRound(1, models.West).Return(3)
				mockGame.EXPECT().PlayerBySeatInRound(1, models.North).Return(4)
				// CurrentWind is called for each player
				mockGame.EXPECT().CurrentWind(1).Return(models.East)
				mockGame.EXPECT().CurrentWind(2).Return(models.South)
				mockGame.EXPECT().CurrentWind(3).Return(models.West)
				mockGame.EXPECT().CurrentWind(4).Return(models.North)
				mockGame.EXPECT().ID().Return(100)
				mockGame.EXPECT().Status().Return(models.GameLive)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(mp *MockPlayerRepo, gameID int) {
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			checkView: func(t *testing.T, view models.GameView) {
				assert.Equal(t, 100, view.ID)
				assert.Equal(t, models.GameLive, view.Status)
				assert.Equal(t, 1, view.Round.Num)
				assert.Len(t, view.PlayerScores, 4)
				assert.Equal(t, 1, view.PlayerScores[0].ID)
				assert.Equal(t, "Alice", view.PlayerScores[0].Name)
				assert.Equal(t, models.East.Kanji(), view.PlayerScores[0].Wind)
				assert.Equal(t, 10, view.PlayerScores[0].Score)
			},
		},
		"game not found": {
			gameID: 999,
			setupGame: func(s *Service, t *testing.T) int {
				// Don't add any games
				return 999
			},
			mockSetup: func(mp *MockPlayerRepo, gameID int) {
				// No player repo calls expected
			},
			expectedErr: assert.Error,
		},
		"player repo error": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Score().Return(models.GameScore{
					1: {GP: 10},
					2: {GP: -5},
					3: {GP: 0},
					4: {GP: -5},
				})
				mockGame.EXPECT().Round().Return(1)
				mockGame.EXPECT().PlayerBySeatInRound(1, models.East).Return(1)
				// PlayerBySeatInRound for other winds won't be called if Get fails
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(mp *MockPlayerRepo, gameID int) {
				mp.On("Get", mock.Anything, 1).Return(models.Player{}, errors.New("player not found"))
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
			testCase.mockSetup(mockPlayerRepo, gameID)

			ctx := context.Background()
			view, err := service.ViewLive(ctx, gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkView != nil {
				testCase.checkView(t, view)
			}
		})
	}
}
