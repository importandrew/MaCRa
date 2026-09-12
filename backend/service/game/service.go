package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type GameRepo interface {
	Status(ctx context.Context, gameID int) (models.GameStatus, error)
	SetStatus(ctx context.Context, gameID int, status models.GameStatus) error
	Players(ctx context.Context, gameID int) (models.StartingPlayers, error)
	GameLog(ctx context.Context, gameID int) ([]models.RoundLog, error)
	SaveRoundOutcome(ctx context.Context, gameID int, data models.RoundData) error
	DeleteLastRound(ctx context.Context, gameID int) error
	Unfinalize(ctx context.Context, gameID int) error
	NewGame(ctx context.Context, seating [4]int, eventID int) (int, error)
	GameResult(ctx context.Context, gameID int) (models.GameResult, error)
	GameResults(ctx context.Context, gameIDs []int) (map[int]models.GameResult, error)
	Finalize(ctx context.Context, data models.FinalGameData) error
	LastGame(ctx context.Context, playerID, eventID int) (int, error)
	CurrentGame(ctx context.Context, playerID, eventID int) (int, error)
	EventID(ctx context.Context, gameID int) (int, error)
	IsPlayerInGame(ctx context.Context, gameID, playerID int) (bool, error)
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
}

func NewService(gameRepo GameRepo, playerRepo PlayerRepo) *Service {
	return &Service{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
	}
}
