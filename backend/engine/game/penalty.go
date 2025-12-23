package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (g *Game) FalseMahjong(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	if value != 30 && value != 60 {
		return fmt.Errorf("false mahjong penalty should be 30 or 60, got: %d", value)
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

	if forceNextRound {
		g.NextRound()
	}

	return nil
}

func (g *Game) Penalty(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	if forceNextRound {
		defer g.NextRound()
	}

	ps := g.score[loser]
	ps.GP -= value
	g.score[loser] = ps

	return nil
}
