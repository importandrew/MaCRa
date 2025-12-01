package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLiveGame(t *testing.T) {
	testData := map[string]struct {
		gameID      int
		setupGame   func(*Service, *testing.T) int
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine)
	}{
		"ok": {
			gameID: 100,
			setupGame: func(s *Service, t *testing.T) int {
				// Create a mock game engine and add it to liveGames
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().ID().Return(100)
				s.liveGames[100] = mockGame
				return 100
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine) {
				assert.NotNil(t, g)
				assert.Equal(t, 100, g.ID())
			},
		},
		"game not found": {
			gameID: 999,
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

			expectedGameID := testCase.setupGame(service, t)

			g, err := service.GetLiveGame(expectedGameID)

			testCase.expectedErr(t, err)
			if err != nil {
				assert.Nil(t, g)
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g)
			}
		})
	}
}

func TestCurrentGame(t *testing.T) {
	testData := map[string]struct {
		playerID    int
		setupGame   func(*Service, *testing.T)
		expectedErr assert.ErrorAssertionFunc
		checkGame   func(*testing.T, GameEngine, int)
	}{
		"ok": {
			playerID: 1,
			setupGame: func(s *Service, t *testing.T) {
				mockGame := NewMockGameEngine(t)
				// CurrentGame calls Players() during iteration, and checkGame calls it again
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4}).Times(2)
				s.liveGames[100] = mockGame
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine, playerID int) {
				assert.NotNil(t, g)
				// Verify player is in the game
				players := g.Players()
				found := false
				for _, id := range players {
					if id == playerID {
						found = true
						break
					}
				}
				assert.True(t, found, "player should be in the game")
			},
		},
		"player not in any game": {
			playerID: 999,
			setupGame: func(s *Service, t *testing.T) {
				mockGame := NewMockGameEngine(t)
				mockGame.EXPECT().Players().Return([4]int{1, 2, 3, 4})
				s.liveGames[100] = mockGame
			},
			expectedErr: assert.Error,
		},
		"no live games": {
			playerID: 1,
			setupGame: func(s *Service, t *testing.T) {
				// Don't add any games
			},
			expectedErr: assert.Error,
		},
		"multiple games": {
			playerID: 5,
			setupGame: func(s *Service, t *testing.T) {
				// First game - Players() may or may not be called depending on map iteration order
				mockGame1 := NewMockGameEngine(t)
				mockGame1.EXPECT().Players().Return([4]int{1, 2, 3, 4}).Maybe()
				s.liveGames[100] = mockGame1

				// Second game with player 5 - Players() will always be called (either first or after game1)
				// If checked first, it finds player 5 and returns early
				// If checked second, it finds player 5 after game1 doesn't have the player
				mockGame2 := NewMockGameEngine(t)
				mockGame2.EXPECT().Players().Return([4]int{5, 6, 7, 8})
				mockGame2.EXPECT().ID().Return(200)
				s.liveGames[200] = mockGame2
			},
			expectedErr: assert.NoError,
			checkGame: func(t *testing.T, g GameEngine, playerID int) {
				assert.NotNil(t, g)
				assert.Equal(t, 200, g.ID())
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.setupGame(service, t)

			g, err := service.CurrentGame(testCase.playerID)

			testCase.expectedErr(t, err)
			if err != nil {
				assert.Nil(t, g)
				return
			}

			if testCase.checkGame != nil {
				testCase.checkGame(t, g, testCase.playerID)
			}
		})
	}
}
