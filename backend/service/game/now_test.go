package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNow(t *testing.T) {
	testData := map[string]struct {
		playerID  int
		setupGame func(*Service, *testing.T) int

		mockSetup func(*MockPlayerRepo, int, int)

		expectedErr assert.ErrorAssertionFunc
		checkNow    func(*testing.T, models.GameNow)
	}{
		"ok": {
			playerID: 1,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() to find the player
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				// Now() calls Left, Opposite, Right, then CurrentWind for each, then Score, ID, Round
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Right(1).Return(2)
				// CurrentWind is called in order: playerID, leftID, oppID, rightID
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
				mockGame.EXPECT().Round().Return(1)
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(mp *MockPlayerRepo, playerID, gameID int) {
				// Player 1 is East, so:
				// Left = 4 (North), Right = 2 (South), Opp = 3 (West)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, 1, now.Round)
				assert.Equal(t, "Alice", now.Me.Name)
				assert.Equal(t, models.East.Kanji(), now.Me.Wind)
			},
		},
		"player not in game": {
			playerID: 999,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() to find the player
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(mp *MockPlayerRepo, playerID, gameID int) {
				// No player repo calls expected
			},
			expectedErr: assert.Error,
		},
		"player repo error": {
			playerID: 1,
			setupGame: func(s *Service, t *testing.T) int {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() to find the player
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				// Now() calls Left, Opposite, Right, then CurrentWind for each
				// If playerRepo.Get fails, CurrentWind won't be called, but we need to set up
				// the calls that happen before the error
				mockGame.EXPECT().Left(1).Return(4)
				mockGame.EXPECT().Opposite(1).Return(3)
				mockGame.EXPECT().Right(1).Return(2)
				// CurrentWind calls happen after playerRepo.Get, so they won't be called if Get fails
				// But Score, ID, Round are called after CurrentWind, so they also won't be called
				s.liveGames[100] = mockGame
				return 100
			},
			mockSetup: func(mp *MockPlayerRepo, playerID, gameID int) {
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
			testCase.mockSetup(mockPlayerRepo, testCase.playerID, gameID)

			ctx := context.Background()
			now, err := service.Now(ctx, testCase.playerID)

			testCase.expectedErr(t, err)
			if err != nil {
				mockPlayerRepo.AssertExpectations(t)
				return
			}

			if testCase.checkNow != nil {
				testCase.checkNow(t, now)
			}
			mockPlayerRepo.AssertExpectations(t)
		})
	}
}
