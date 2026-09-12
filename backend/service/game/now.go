package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Now(ctx context.Context, playerID, eventID int) (models.GameNow, error) {
	game, err := s.CurrentGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}

	leftID := game.Left(playerID)
	oppID := game.Opposite(playerID)
	rightID := game.Right(playerID)

	me, err := s.playerRepo.Get(ctx, playerID)
	if err != nil {
		return models.GameNow{}, err
	}
	left, err := s.playerRepo.Get(ctx, leftID)
	if err != nil {
		return models.GameNow{}, err
	}
	right, err := s.playerRepo.Get(ctx, rightID)
	if err != nil {
		return models.GameNow{}, err
	}
	opp, err := s.playerRepo.Get(ctx, oppID)
	if err != nil {
		return models.GameNow{}, err
	}

	myWind := game.CurrentWind(playerID)
	leftWind := game.CurrentWind(leftID)
	oppWind := game.CurrentWind(oppID)
	rightWind := game.CurrentWind(rightID)

	score := game.Score()

	myScore := score[playerID].GP
	rightScore := score[rightID].GP
	oppScore := score[oppID].GP
	leftScore := score[leftID].GP

	round := game.Round()

	now := models.GameNow{
		ID:     game.ID(),
		Status: game.Status(),
		Round: models.Round{
			Num:   round,
			Kanji: models.RoundKanji(round),
		},
		Me: models.PlayerWindScore{
			ID:    playerID,
			Name:  me.Name,
			Wind:  myWind.Kanji(),
			Score: myScore,
		},
		Right: models.PlayerWindScore{
			ID:    rightID,
			Name:  right.Name,
			Wind:  rightWind.Kanji(),
			Score: rightScore,
		},
		Opp: models.PlayerWindScore{
			ID:    oppID,
			Name:  opp.Name,
			Wind:  oppWind.Kanji(),
			Score: oppScore,
		},
		Left: models.PlayerWindScore{
			ID:    leftID,
			Name:  left.Name,
			Wind:  leftWind.Kanji(),
			Score: leftScore,
		},
	}

	return now, nil
}
