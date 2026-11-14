package game

import (
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

func (g *Game) FalseMahjong(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	if !g.isSeated(loser) {
		return fmt.Errorf("%w: false mahjong: player %d is not in the game", models.ErrValidatePenalty, loser)
	}

	if !slices.Contains(models.FalseMahjongAmounts, value) {
		return fmt.Errorf("%w: false mahjong penalty should be one of %v, got: %d",
			models.ErrValidatePenalty, models.FalseMahjongAmounts, value)
	}

	for id, ps := range g.score {
		switch id {
		case loser:
			ps.GP -= value
			g.score[id] = ps
		default:
			ps.GP += value / 3
			g.score[id] = ps
		}
	}

	g.logSize++

	if forceNextRound {
		return g.NextRound()
	}

	return nil
}

func (g *Game) Penalty(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	if !g.isSeated(loser) {
		return fmt.Errorf("%w: penalty: player %d is not in the game", models.ErrValidatePenalty, loser)
	}

	if value <= 0 {
		return fmt.Errorf("%w: penalty should be positive, got: %d", models.ErrValidatePenalty, value)
	}

	ps := g.score[loser]
	ps.GP -= value
	g.score[loser] = ps

	g.logSize++

	if forceNextRound {
		return g.NextRound()
	}

	return nil
}
