package game

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test full game progression.
func TestSeating_FullGame(t *testing.T) {
	g, _ := NewGame(Config{
		ID:      100,
		Players: [4]int{1, 2, 3, 4},
	})

	expectedOrder := [][]int{
		// East
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{3, 4, 1, 2},
		{4, 1, 2, 3},

		// South
		{2, 1, 4, 3},
		{1, 4, 3, 2},
		{4, 3, 2, 1},
		{3, 2, 1, 4},

		// West
		{3, 4, 2, 1},
		{4, 2, 1, 3},
		{2, 1, 3, 4},
		{1, 3, 4, 2},

		// North
		{4, 3, 1, 2},
		{3, 1, 2, 4},
		{1, 2, 4, 3},
		{2, 4, 3, 1},
	}

	for i, order := range expectedOrder {
		playersInWindOrder := make([]int, 0, 4)
		for _, id := range g.currentSeating {
			playersInWindOrder = append(playersInWindOrder, id)
		}

		require.True(t, slices.Equal(order, playersInWindOrder), "Round %d: expected order %v, got %v", i+1, order, playersInWindOrder)

		g.NextRound()
	}
}
