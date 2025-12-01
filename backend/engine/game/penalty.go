package game

import (
	"errors"
	"fmt"
)

func (g *Game) FalseMahjong(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return errors.New("game is not live")
	}

	if value != 30 && value != 60 {
		return fmt.Errorf("false mahjong penalty should be 30 or 60, got: %d", value)
	}

	if forceNextRound {
		defer g.NextRound()
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

	return nil
}

func (g *Game) Penalty(loser, value int, forceNextRound bool) error {
	if !g.IsLive() {
		return errors.New("game is not live")
	}

	if forceNextRound {
		defer g.NextRound()
	}

	ps := g.score[loser]
	ps.GP -= value
	g.score[loser] = ps

	return nil
}
