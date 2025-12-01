package game

import (
	"github.com/importandrew/MaCRa/models"
)

// NewScore initialises Score with players.
func NewScore(seating Seating) models.GameScore {
	score := make(models.GameScore, 4)
	for _, id := range seating {
		score[id] = models.PlayerScore{}
	}
	return score
}
