package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// Now returns the player's table: their live game, or a tournament game
// whose tour is published but not started (status created; no rounds can be
// entered yet), with its tour and table.
func (s *Service) Now(ctx context.Context, playerID, eventID int) (models.GameNow, error) {
	table, err := s.seatedTable(ctx, playerID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}
	game, log, status, tour, limit := table.game, table.log, table.status, table.tour, table.limit

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

	names := map[int]string{playerID: me.Name, leftID: left.Name, rightID: right.Name, oppID: opp.Name}

	undone, err := s.gameRepo.UndoneEntries(ctx, game.ID())
	if err != nil {
		return models.GameNow{}, err
	}

	now := models.GameNow{
		ID:      game.ID(),
		Status:  game.Status(),
		LogSize: game.LogSize(),
		Rules:   game.Rules(),
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
		Log:   gameLog(game, log, undone, names),
		Tour:  tour,
		Timer: gameTimer(game, limit),
	}
	// The stored status: a game stopped at its tour's time has ended in the
	// replay but stays live until it is finished (the timer says so).
	now.Status = status

	return now, nil
}

// seatedTableState is what Now shows besides the seats: the player's game,
// its stored status, its tour seat and time limit.
type seatedTableState struct {
	game   GameEngine
	log    []models.RoundLog
	status models.GameStatus
	tour   *models.TourSeat
	limit  *models.GameTimeLimit
}

// seatedTable loads the player's game that isn't over (see SeatedGame) with
// its stored status, tour seat and time limit.
func (s *Service) seatedTable(ctx context.Context, playerID, eventID int) (seatedTableState, error) {
	gameID, err := s.gameRepo.SeatedGame(ctx, playerID, eventID)
	if err != nil {
		return seatedTableState{}, err
	}
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return seatedTableState{}, err
	}
	game, log, err := s.load(ctx, gameID)
	if err != nil {
		return seatedTableState{}, fmt.Errorf("load game %d: %w", gameID, err)
	}
	tour, err := s.gameRepo.TourSeat(ctx, gameID)
	if err != nil {
		return seatedTableState{}, err
	}
	limit, err := s.gameRepo.TimeLimit(ctx, gameID)
	if err != nil {
		return seatedTableState{}, err
	}
	return seatedTableState{game: game, log: log, status: status, tour: tour, limit: limit}, nil
}
