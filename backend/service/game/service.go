package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type GameRepo interface {
	Status(ctx context.Context, gameID int) (models.GameStatus, error)
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
	Status() models.GameStatus
	IsLive() bool
	Round() int
	End() error
	NextRound() error
	Score() models.GameScore

	EnterRoundResult(outcome models.Outcome, hand int, loser int, winner int) error
	FalseMahjong(loser int, value int, forceNextRound bool) error
	Penalty(loser int, value int, forceNextRound bool) error

	Left(id int) int
	Right(id int) int
	Opposite(id int) int

	StartWind(id int) models.Wind
	CurrentWind(id int) models.Wind
	RoundStartWind(id int) models.Wind

	Players() [4]int
	RoundStartPlayers() [4]int

	PlayerBySeatInRound(round int, seat models.Wind) int
	PlayerByStartingSeat(seat models.Wind) int
	PlayerByWindInRound(round int, wind models.Wind) int
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
