package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) FalseMahjong(ctx context.Context, gameID, playerID, value int, forceNextRound bool) error {
	game, err := s.GetLiveGame(gameID)
	if err != nil {
		return err
	}

	err = game.FalseMahjong(playerID, value, forceNextRound)
	if err != nil {
		return err
	}

	players := game.Players()
	score := game.Score()

	roundOutcome := models.RoundData{
		Round:           game.Round(),
		Outcome:         models.FalseMahjong,
		LoserID:         playerID,
		Value:           value,
		EastStartScore:  score[players[models.East]].GP,
		SouthStartScore: score[players[models.South]].GP,
		WestStartScore:  score[players[models.West]].GP,
		NorthStartScore: score[players[models.North]].GP,
		ForceNextRound:  forceNextRound,
	}

	err = s.gameRepo.SaveRoundOutcome(ctx, game.ID(), roundOutcome)
	if err != nil {
		return err
	}

	return nil
}
