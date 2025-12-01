package game

import (
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func TestGame_NextRound(t *testing.T) {
	testData := map[string]struct {
		setupGame   func() *Game
		expectedErr assert.ErrorAssertionFunc
		check       func(*testing.T, *Game)
	}{
		"ok": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, 2, g.Round())
				assert.Equal(t, models.East, g.CurrentWind(2))
				assert.Equal(t, models.South, g.CurrentWind(3))
				assert.Equal(t, models.West, g.CurrentWind(4))
				assert.Equal(t, models.North, g.CurrentWind(1))
			},
		},
		"ok round 5": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				for range 3 {
					g.NextRound()
				}
				return g
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, 5, g.Round())
				assert.Equal(t, models.South, g.CurrentWind(1))
				assert.Equal(t, models.East, g.CurrentWind(2))
				assert.Equal(t, models.North, g.CurrentWind(3))
				assert.Equal(t, models.West, g.CurrentWind(4))
			},
		},
		"ok round 9": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				for range 7 {
					g.NextRound()
				}
				return g
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, 9, g.Round())
				assert.Equal(t, models.North, g.CurrentWind(1))
				assert.Equal(t, models.West, g.CurrentWind(2))
				assert.Equal(t, models.East, g.CurrentWind(3))
				assert.Equal(t, models.South, g.CurrentWind(4))
			},
		},
		"ok round 12": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				for range 11 {
					g.NextRound()
				}
				return g
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, 13, g.Round())
			},
		},
		"ok end": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				for range 15 {
					g.NextRound()
				}
				return g
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, models.GameEnded, g.Status())
				assert.False(t, g.IsLive())
			},
		},
		"game not live": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				g.End()
				return g
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.NextRound()
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}
			if testCase.check != nil {
				testCase.check(t, g)
			}
		})
	}
}

func TestGame_EnterRoundResult(t *testing.T) {
	testData := map[string]struct {
		setupGame   func() *Game
		hand        int
		loser       int
		winner      int
		expectedErr assert.ErrorAssertionFunc
		check       func(*testing.T, *Game)
	}{
		"ok deal in": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			hand:        8,
			loser:       2,
			winner:      1,
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 32, score[1].GP)
				assert.Equal(t, -16, score[2].GP)
				assert.Equal(t, -8, score[3].GP)
				assert.Equal(t, -8, score[4].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"ok self": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			hand:        10,
			loser:       0,
			winner:      1,
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 54, score[1].GP)
				assert.Equal(t, -18, score[2].GP)
				assert.Equal(t, -18, score[3].GP)
				assert.Equal(t, -18, score[4].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"ok draw": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			hand:        0,
			loser:       0,
			winner:      0,
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 0, score[1].GP)
				assert.Equal(t, 0, score[2].GP)
				assert.Equal(t, 0, score[3].GP)
				assert.Equal(t, 0, score[4].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"invalid hand deal in": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			hand:        5,
			loser:       2,
			winner:      1,
			expectedErr: assert.Error,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 0, score[1].GP)
				assert.Equal(t, 0, score[2].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"invalid hand self": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			hand:        7,
			loser:       0,
			winner:      1,
			expectedErr: assert.Error,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 0, score[1].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"game not live": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				g.End()
				return g
			},
			hand:        8,
			loser:       2,
			winner:      1,
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.EnterRoundResult(testCase.hand, testCase.loser, testCase.winner)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}
			if testCase.check != nil {
				testCase.check(t, g)
			}
		})
	}
}
