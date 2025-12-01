package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_Penalty(t *testing.T) {
	testData := map[string]struct {
		setupGame      func() *Game
		loser          int
		value          int
		forceNextRound bool
		expectedErr    assert.ErrorAssertionFunc
		check          func(*testing.T, *Game)
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
			loser:          2,
			value:          10,
			forceNextRound: false,
			expectedErr:    assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, -10, score[2].GP)
				assert.Equal(t, 0, score[1].GP)
				assert.Equal(t, 0, score[3].GP)
				assert.Equal(t, 0, score[4].GP)
				assert.Equal(t, 1, g.Round())
			},
		},
		"ok with force next round": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			loser:          3,
			value:          20,
			forceNextRound: true,
			expectedErr:    assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, -20, score[3].GP)
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
			loser:          2,
			value:          10,
			forceNextRound: false,
			expectedErr:    assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.Penalty(testCase.loser, testCase.value, testCase.forceNextRound)
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

func TestGame_FalseMahjong(t *testing.T) {
	testData := map[string]struct {
		setupGame      func() *Game
		loser          int
		value          int
		forceNextRound bool
		expectedErr    assert.ErrorAssertionFunc
		check          func(*testing.T, *Game)
	}{
		"ok value 30": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			loser:          2,
			value:          30,
			forceNextRound: false,
			expectedErr:    assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, -30, score[2].GP)
				assert.Equal(t, 10, score[1].GP)
				assert.Equal(t, 10, score[3].GP)
				assert.Equal(t, 10, score[4].GP)
				assert.Equal(t, 1, g.Round())
			},
		},
		"ok value 60 with force next round": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			loser:          3,
			value:          60,
			forceNextRound: true,
			expectedErr:    assert.NoError,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, -60, score[3].GP)
				assert.Equal(t, 20, score[1].GP)
				assert.Equal(t, 20, score[2].GP)
				assert.Equal(t, 20, score[4].GP)
				assert.Equal(t, 2, g.Round())
			},
		},
		"invalid value": {
			setupGame: func() *Game {
				g, _ := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   Rules{},
				})
				return g
			},
			loser:          2,
			value:          50,
			forceNextRound: false,
			expectedErr:    assert.Error,
			check: func(t *testing.T, g *Game) {
				score := g.Score()
				assert.Equal(t, 0, score[2].GP)
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
			loser:          2,
			value:          30,
			forceNextRound: false,
			expectedErr:    assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.FalseMahjong(testCase.loser, testCase.value, testCase.forceNextRound)
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
