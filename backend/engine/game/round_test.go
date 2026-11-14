package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestGame_NextRound(t *testing.T) {
	testData := map[string]struct {
		setupGame   func() *Game
		expectedErr assert.ErrorAssertionFunc
		check       func(*testing.T, *Game)
	}{
		"ok": {
			setupGame: func() *Game {
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
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
		outcome     models.Outcome
		hand        int
		loser       int
		winner      int
		expectedErr assert.ErrorAssertionFunc
		check       func(*testing.T, *Game)
	}{
		"ok deal in": {
			setupGame: func() *Game {
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				return g
			},
			outcome:     models.Deal,
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				return g
			},
			outcome:     models.Self,
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				return g
			},
			outcome:     models.Draw,
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				return g
			},
			outcome:     models.Deal,
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				return g
			},
			outcome:     models.Self,
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
				g := NewGame(Config{
					ID:      100,
					Players: [4]int{1, 2, 3, 4},
					Rules:   models.StandardGameRules(),
				})
				g.End()
				return g
			},
			outcome:     models.Deal,
			hand:        8,
			loser:       2,
			winner:      1,
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.EnterRoundResult(testCase.outcome, testCase.hand, testCase.loser, testCase.winner)
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

func TestGame_EnterRoundResult_InvalidPlayers(t *testing.T) {
	testData := map[string]struct {
		outcome models.Outcome
		loser   int
		winner  int
	}{
		"self with an unseated winner":    {outcome: models.Self, winner: 99},
		"deal with an unseated winner":    {outcome: models.Deal, loser: 2, winner: 99},
		"deal with an unseated loser":     {outcome: models.Deal, loser: 99, winner: 1},
		"deal with winner equal to loser": {outcome: models.Deal, loser: 1, winner: 1},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := NewGame(Config{ID: 100, Players: [4]int{1, 2, 3, 4}})

			err := g.EnterRoundResult(testCase.outcome, 8, testCase.loser, testCase.winner)

			assert.ErrorIs(t, err, models.ErrValidateEnterRoundResult)
			assert.Equal(t, NewScore([4]int{1, 2, 3, 4}), g.Score())
			assert.Equal(t, 1, g.Round())
		})
	}
}

// TestGame_ByRules checks that the game follows its rules: the number of
// hands and the minimum hand.
func TestGame_ByRules(t *testing.T) {
	newGame := func(rules models.GameRules) *Game {
		return NewGame(Config{ID: 100, Players: [4]int{1, 2, 3, 4}, Rules: rules})
	}

	for _, hands := range models.ValidHands {
		t.Run(fmt.Sprintf("ends after %d hands", hands), func(t *testing.T) {
			g := newGame(models.GameRules{Hands: hands, MinHand: 8})
			for range hands - 1 {
				assert.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
			}
			assert.True(t, g.IsLive())
			assert.Equal(t, hands, g.Round())

			assert.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
			assert.Equal(t, models.GameEnded, g.Status())
		})
	}

	t.Run("minimum hand", func(t *testing.T) {
		g := newGame(models.GameRules{Hands: 16, MinHand: 24})
		assert.ErrorIs(t, g.EnterRoundResult(models.Self, 23, 0, 1), models.ErrValidateEnterRoundResult)
		assert.ErrorIs(t, g.EnterRoundResult(models.Deal, 23, 2, 1), models.ErrValidateEnterRoundResult)
		assert.NoError(t, g.EnterRoundResult(models.Deal, 24, 2, 1))
	})
}
