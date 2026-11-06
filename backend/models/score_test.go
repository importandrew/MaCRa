package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameScore_CalculateTP(t *testing.T) {
	testData := map[string]struct {
		gp       [4]int // game points of players 1..4
		expected [4]float32
	}{
		"4 < 3 < 2 < 1": {gp: [4]int{30, 10, -10, -30}, expected: [4]float32{4, 2, 1, 0}},
		"4 = 3 < 2 < 1": {gp: [4]int{30, 10, -20, -20}, expected: [4]float32{4, 2, 0.5, 0.5}},
		"4 < 3 = 2 < 1": {gp: [4]int{30, 0, 0, -30}, expected: [4]float32{4, 1.5, 1.5, 0}},
		"4 < 3 < 2 = 1": {gp: [4]int{20, 20, -10, -30}, expected: [4]float32{3, 3, 1, 0}},
		"4 = 3 = 2 < 1": {gp: [4]int{30, -10, -10, -10}, expected: [4]float32{4, 1, 1, 1}},
		"4 < 3 = 2 = 1": {gp: [4]int{10, 10, 10, -30}, expected: [4]float32{7.0 / 3, 7.0 / 3, 7.0 / 3, 0}},
		"4 = 3 < 2 = 1": {gp: [4]int{20, 20, -20, -20}, expected: [4]float32{3, 3, 0.5, 0.5}},
		"4 = 3 = 2 = 1": {gp: [4]int{0, 0, 0, 0}, expected: [4]float32{1.75, 1.75, 1.75, 1.75}},
		"unsorted ids":  {gp: [4]int{-30, 30, -10, 10}, expected: [4]float32{0, 4, 1, 2}},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			score := GameScore{}
			for i, gp := range testCase.gp {
				score[i+1] = PlayerScore{GP: gp}
			}

			score.CalculateTP()

			for i, tp := range testCase.expected {
				assert.InDelta(t, tp, score[i+1].TP, 1e-6, "player %d", i+1)
				assert.Equal(t, testCase.gp[i], score[i+1].GP, "player %d GP changed", i+1)
			}
		})
	}
}
