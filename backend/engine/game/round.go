package game

import (
	"errors"
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
		g.roundSeating.ApplySouthNorthSeating()
		g.currentSeating = g.roundSeating
	case 8:
		g.roundSeating.ApplyWestSeating()
		g.currentSeating = g.roundSeating
	case 12:
		g.roundSeating.ApplySouthNorthSeating()
		g.currentSeating = g.roundSeating
	case 16:
		return g.End()
	default:
		g.currentSeating.NextRoundSeating()
	}
	g.round += 1
	logger.WithField("next", g.round).Info("game moved to the next round")

	return nil
}

func (g *Game) EnterRoundResult(outcome models.Outcome, hand, loser, winner int) error {
	if !g.IsLive() {
		return models.ErrGameNotLive
	}

	switch outcome {
	case models.Draw:
		if err := validateDraw(hand, loser, winner); err != nil {
			return fmt.Errorf("%w: draw: %w", models.ErrValidateEnterRoundResult, err)
		}
		g.draw()
	case models.Self:
		if err := validateSelf(hand, loser, winner); err != nil {
			return fmt.Errorf("%w: self: %w", models.ErrValidateEnterRoundResult, err)
		}
		g.self(hand, winner)
	case models.Deal:
		if err := validateDeal(hand, loser, winner); err != nil {
			return fmt.Errorf("%w: deal: %w", models.ErrValidateEnterRoundResult, err)
		}
		g.deal(hand, loser, winner)
	default:
		return fmt.Errorf("%w: invalid outcome: %s", models.ErrValidateEnterRoundResult, outcome)
	}

	err := g.NextRound()
	if err != nil {
		return err
	}

	return nil
}

func validateDraw(hand, loser, winner int) error {
	if hand != 0 {
		return errors.New("hand score should be zero")
	}

	if winner != 0 {
		return errors.New("absent winner is required")
	}

	if loser != 0 {
		return errors.New("absent loser is required")
	}

	return nil
}

func (g *Game) draw() {
	logger.WithField("result", "draw").Info(g.score.String())
}

func validateSelf(hand, loser, winner int) error {
	if hand < 8 {
		return errors.New("hand score should be 8 or higher")
	}

	if winner == 0 {
		return errors.New("winner is required")
	}

	if loser != 0 {
		return errors.New("absent loser is required")
	}

	return nil
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

func validateDeal(hand, loser, winner int) error {
	if hand < 8 {
		return errors.New("hand score should be 8 or higher")
	}

	if winner == 0 {
		return errors.New("winner is required")
	}

	if loser == 0 {
		return errors.New("loser is required")
	}

	return nil
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
