package game

import (
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func TestGame_Players(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	players := g.Players()
	assert.Equal(t, [4]int{1, 2, 3, 4}, players)

	// Players should not change even after rounds
	g.NextRound()
	players = g.Players()
	assert.Equal(t, [4]int{1, 2, 3, 4}, players)
}

func TestGame_Left(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	// Player 1 is East, Left is North (player 4)
	assert.Equal(t, 4, g.Left(1))
	// Player 2 is South, Left is East (player 1)
	assert.Equal(t, 1, g.Left(2))
	// Player 3 is West, Left is South (player 2)
	assert.Equal(t, 2, g.Left(3))
	// Player 4 is North, Left is West (player 3)
	assert.Equal(t, 3, g.Left(4))
}

func TestGame_Right(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	// Player 1 is East, Right is South (player 2)
	assert.Equal(t, 2, g.Right(1))
	// Player 2 is South, Right is West (player 3)
	assert.Equal(t, 3, g.Right(2))
	// Player 3 is West, Right is North (player 4)
	assert.Equal(t, 4, g.Right(3))
	// Player 4 is North, Right is East (player 1)
	assert.Equal(t, 1, g.Right(4))
}

func TestGame_Opposite(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	// Player 1 is East, Opposite is West (player 3)
	assert.Equal(t, 3, g.Opposite(1))
	// Player 2 is South, Opposite is North (player 4)
	assert.Equal(t, 4, g.Opposite(2))
	// Player 3 is West, Opposite is East (player 1)
	assert.Equal(t, 1, g.Opposite(3))
	// Player 4 is North, Opposite is South (player 2)
	assert.Equal(t, 2, g.Opposite(4))
}

func TestGame_CurrentWind(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	// Initial seating: 1=East, 2=South, 3=West, 4=North
	assert.Equal(t, models.East, g.CurrentWind(1))
	assert.Equal(t, models.South, g.CurrentWind(2))
	assert.Equal(t, models.West, g.CurrentWind(3))
	assert.Equal(t, models.North, g.CurrentWind(4))

	// After one round, seating rotates
	g.NextRound()
	// New seating: 2=East, 3=South, 4=West, 1=North
	assert.Equal(t, models.East, g.CurrentWind(2))
	assert.Equal(t, models.South, g.CurrentWind(3))
	assert.Equal(t, models.West, g.CurrentWind(4))
	assert.Equal(t, models.North, g.CurrentWind(1))
}

func TestGame_StartWind(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
		Rules:   Rules{},
	})

	// Start wind should always be the same regardless of rounds
	assert.Equal(t, models.East, g.StartWind(1))
	assert.Equal(t, models.South, g.StartWind(2))
	assert.Equal(t, models.West, g.StartWind(3))
	assert.Equal(t, models.North, g.StartWind(4))

	// After rounds, start wind should not change
	g.NextRound()
	assert.Equal(t, models.East, g.StartWind(1))
	assert.Equal(t, models.South, g.StartWind(2))
	assert.Equal(t, models.West, g.StartWind(3))
	assert.Equal(t, models.North, g.StartWind(4))
}
