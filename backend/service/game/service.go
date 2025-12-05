package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type GameRepo interface {
	SetStatus(ctx context.Context, gameID int, status models.GameStatus) error
	LiveGames(ctx context.Context) ([]int, error)
	Players(ctx context.Context, gameID int) (models.StartingPlayers, error)
	GameLog(ctx context.Context, gameID int) ([]models.RoundLog, error)
	SaveRoundOutcome(ctx context.Context, gameID int, data models.RoundData) error
	NewGame(ctx context.Context, seating [4]int) (int, error)
	GameResult(ctx context.Context, gameID int) (models.GameResult, error)
	Finalize(ctx context.Context, data models.FinalGameData) error
	LastGame(ctx context.Context, playerID int) (int, error)
}

type PlayerRepo interface {
	Get(ctx context.Context, id int) (models.Player, error)
}

type GameEngine interface {
	ID() int
	Round() int
	Status() models.GameStatus
	IsLive() bool
	CurrentWind(id int) models.Wind
	End() error
	EnterRoundResult(hand int, loser int, winner int) error
	FalseMahjong(loser int, value int, forceNextRound bool) error
	Left(id int) int
	NextRound() error
	Opposite(id int) int
	Penalty(loser int, value int, forceNextRound bool) error
	Players() [4]int
	Right(id int) int
	StartWind(id int) models.Wind
	Score() models.GameScore
}

type Service struct {
	gameRepo   GameRepo
	playerRepo PlayerRepo
	liveGames  map[int]GameEngine
}

func NewService(gameRepo GameRepo, playerRepo PlayerRepo) *Service {
	return &Service{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
		liveGames:  make(map[int]GameEngine),
	}
}
