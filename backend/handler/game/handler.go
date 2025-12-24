package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type GameService interface {
	FalseMahjong(ctx context.Context, playerID, penalizedID, value int, forceNextRound bool) (models.GameNow, error)
	Penalty(ctx context.Context, playerID, penalizedID, value int, forceNextRound bool) (models.GameNow, error)
	EnterRound(ctx context.Context, playerID int, outcome models.Outcome, value, winnerID, loserID int) (models.GameNow, error)
	Now(ctx context.Context, playerID int) (models.GameNow, error)
	LastGameResult(ctx context.Context, playerID int) (models.GameResult, error)

	Abort(ctx context.Context, gameID int) error
	GameResult(ctx context.Context, gameID int) (models.GameResult, error)
	Logs(ctx context.Context, gameID int) (models.GameLog, error)
	NewGame(ctx context.Context, playerIDs []int) (gameID int, err error)
	Players(ctx context.Context, gameID int) (models.GamePlayers, error)
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
