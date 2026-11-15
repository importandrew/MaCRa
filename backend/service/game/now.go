package game

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// Now returns the user's table: the live game they sit at (their own seat or
// one they took as a substitute), or a tournament game whose tour is
// published but not started (status created; no rounds can be entered yet),
// with its tour and table. Seats keep their players' IDs; the names are of
// whoever sits there.
func (s *Service) Now(ctx context.Context, userID, eventID int) (models.GameNow, error) {
	table, err := s.seatedTable(ctx, userID, eventID)
	if err != nil {
		return models.GameNow{}, err
	}
	game, log, status, tour, limit := table.game, table.log, table.status, table.tour, table.limit

	// The table is laid out from the seat the user sits at: their own, or
	// the one they took as a substitute.
	seat, err := s.gameRepo.SeatOf(ctx, game.ID(), userID)
	if err != nil {
		return models.GameNow{}, err
	}
	leftID := game.Left(seat)
	oppID := game.Opposite(seat)
	rightID := game.Right(seat)

	names, err := s.seatNames(ctx, game.ID())
	if err != nil {
		return models.GameNow{}, err
	}

	score := game.Score()
	round := game.Round()
	seatScore := func(id int) models.PlayerWindScore {
		return models.PlayerWindScore{ID: id, Name: names[id], Wind: game.CurrentWind(id).Kanji(), Score: score[id].GP}
	}

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
		Me:    seatScore(seat),
		Right: seatScore(rightID),
		Opp:   seatScore(oppID),
		Left:  seatScore(leftID),
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

// seatNames returns the name of whoever sits at each seat of the game now,
// by the seat's player ID: the player, or the substitute who took the seat.
func (s *Service) seatNames(ctx context.Context, gameID int) (map[int]string, error) {
	sitters, err := s.gameRepo.Sitters(ctx, gameID)
	if err != nil {
		return nil, err
	}
	names := make(map[int]string, len(sitters))
	// In seat order, so the lookups don't depend on the map's.
	for _, seat := range slices.Sorted(maps.Keys(sitters)) {
		sitter := sitters[seat]
		player, err := s.playerRepo.Get(ctx, sitter)
		if err != nil {
			return nil, err
		}
		names[seat] = player.Name
	}
	return names, nil
}

// SeatSitter seats a substitute at a seat of a tournament game, or the
// seat's player back (see repo SeatSitter).
func (s *Service) SeatSitter(ctx context.Context, seat models.SeatChange) error {
	return s.gameRepo.SeatSitter(ctx, seat)
}

// HasPlace reports whether the user has a place in the game: their own seat,
// or one they sat at as a substitute.
func (s *Service) HasPlace(ctx context.Context, gameID, userID int) (bool, error) {
	return s.gameRepo.HasPlace(ctx, gameID, userID)
}
