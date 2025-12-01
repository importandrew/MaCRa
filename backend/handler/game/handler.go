package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type GameService interface {
	Abort(ctx context.Context, gameID int) error
	GameResult(ctx context.Context, gameID int) (models.GameResult, error)
	Logs(ctx context.Context, gameID int) (models.GameLog, error)
	Now(ctx context.Context, playerID int) (models.GameNow, error)
	NewGame(ctx context.Context, playerIDs []int) (gameID int, err error)
	FalseMahjong(ctx context.Context, gameID, playerID, value int, forceNextRound bool) error
	Penalty(ctx context.Context, gameID, playerID, value int, forceNextRound bool) error
	Players(ctx context.Context, gameID int) (models.GamePlayers, error)
	EnterRound(ctx context.Context, playerID int, outcome models.Outcome, value, winnerID, loserID int) error
}

type PlayerService interface {
	Get(ctx context.Context, id int) (models.Player, error)
}

type Handler struct {
	game   GameService
	player PlayerService
}

func NewHandler(game GameService, player PlayerService) Handler {
	return Handler{
		game:   game,
		player: player,
	}
}
