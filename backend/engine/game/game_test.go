package game

import (
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	testData := map[string]struct {
		config      Config
		expectedErr assert.ErrorAssertionFunc
		check       func(*testing.T, *Game)
	}{
		"ok": {
			config: Config{
				ID:      100,
				Players: [4]int{1, 2, 3, 4},
				Rules:   Rules{},
			},
			expectedErr: assert.NoError,
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, 100, g.ID())
				assert.Equal(t, models.GameLive, g.Status())
				assert.True(t, g.IsLive())
				assert.Equal(t, 1, g.Round())
				assert.Equal(t, [4]int{1, 2, 3, 4}, g.Players())
				score := g.Score()
				assert.Len(t, score, 4)
				for _, id := range [4]int{1, 2, 3, 4} {
					assert.Contains(t, score, id)
					assert.Equal(t, 0, score[id].GP)
					assert.Equal(t, float32(0), score[id].TP)
				}
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			got, err := NewGame(testCase.config)
			testCase.expectedErr(t, err)
			if err != nil {
				assert.Nil(t, got)
				return
			}
			assert.NotNil(t, got)
			if testCase.check != nil {
				testCase.check(t, got)
			}
		})
	}
}

func TestGame_End(t *testing.T) {
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
				assert.Equal(t, models.GameEnded, g.Status())
				assert.False(t, g.IsLive())
			},
		},
		"game already ended": {
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
			check: func(t *testing.T, g *Game) {
				assert.Equal(t, models.GameEnded, g.Status())
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			g := testCase.setupGame()
			err := g.End()
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

func TestGame_IsLive(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	assert.True(t, g.IsLive())

	g.End()
	assert.False(t, g.IsLive())
}

func TestGame_ID(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      123,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	assert.Equal(t, 123, g.ID())
}

func TestGame_Status(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	assert.Equal(t, models.GameLive, g.Status())

	g.End()
	assert.Equal(t, models.GameEnded, g.Status())
}

func TestGame_Round(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	assert.Equal(t, 1, g.Round())

	g.NextRound()
	assert.Equal(t, 2, g.Round())
}

func TestGame_Score(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	score := g.Score()
	assert.NotNil(t, score)
	assert.Len(t, score, 4)
}
