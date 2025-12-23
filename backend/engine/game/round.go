package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

func (g *Game) NextRound() error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}
	switch g.round {
	case 4:
		g.currentSeating = g.startSeating
		g.currentSeating.ApplySouthNorthSeating()
	case 8:
		g.currentSeating = g.startSeating
		g.currentSeating.ApplySouthNorthSeating()
		g.currentSeating.ApplyWestSeating()
	case 12:
		g.currentSeating = g.startSeating
		g.currentSeating.ApplySouthNorthSeating()
		g.currentSeating.ApplyWestSeating()
		g.currentSeating.ApplySouthNorthSeating()
	case 16:
		return g.End()
	default:
		g.currentSeating.NextRoundSeating()
	}
	g.round += 1
	logger.WithField("next", g.round).Info("game moved to the next round")

	return nil
}

func (g *Game) EnterRoundResult(hand, loser, winner int) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	switch {
	case loser == 0 && winner == 0:
		g.draw()
	case loser == 0:
		if hand < 8 {
			return models.ErrMinScore
		}
		g.self(hand, winner)
	default:
		if hand < 8 {
			return models.ErrMinScore
		}
		g.deal(hand, loser, winner)
	}

	err := g.NextRound()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) draw() {
	logger.WithField("result", "draw").Info(g.score.String())
}

func (g *Game) self(hand, winner int) {
	for id, ps := range g.score {
		switch id {
		case winner:
			ps.GP += (hand + 8) * 3
			g.score[id] = ps
		default:
			ps.GP -= (hand + 8)
			g.score[id] = ps
		}
	}

	logger.WithField("result", fmt.Sprintf("self by %d (hand: %d)", winner, hand)).
		Info(g.score.String())
}

func (g *Game) deal(hand, loser, winner int) {
	for id, ps := range g.score {
		switch id {
		case winner:
			ps.GP += hand + 8*3
			g.score[id] = ps
		case loser:
			ps.GP -= (hand + 8)
			g.score[id] = ps
		default:
			ps.GP -= 8
			g.score[id] = ps
		}
	}

	logger.WithField("result", fmt.Sprintf("deal-in by %d into %d (hand: %d)", loser, winner, hand)).
		Info(g.score.String())
}
