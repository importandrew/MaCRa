package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/game"
	"github.com/importandrew/MaCRa/service/player"
)

type Service struct {
	repo          EventRepo
	gameService   *game.Service
	playerService *player.Service
}

type EventRepo interface {
	Create(ctx context.Context, event models.Event) (int, error)
	Get(ctx context.Context, eventID int) (models.Event, error)
	List(ctx context.Context) ([]models.Event, error)
	Update(ctx context.Context, event models.Event) error
	Finish(ctx context.Context, eventID int) error
	ListGames(ctx context.Context, eventID int) ([]models.Game, error)
	AddMemberRoles(ctx context.Context, eventID int, userID int, roles []models.EventRole) error
	RemoveMembers(ctx context.Context, eventID int, userIDs []int) error
	GetMemberRoles(ctx context.Context, eventID int, userID int) ([]models.EventRole, error)
	ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error)
	GameHistory(ctx context.Context, eventID int) ([]int, error)
	PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error)
	ListGamesByStatus(ctx context.Context, eventID int, status models.GameStatus, limit, offset int) ([]models.EventGameItemRow, error)
	ListEndedPlayerResults(ctx context.Context, eventID int) ([]models.EventPlayerResultRow, error)
}

func NewService(repo EventRepo, gameService *game.Service, playerService *player.Service) *Service {
	return &Service{
		repo:          repo,
		gameService:   gameService,
		playerService: playerService,
	}
}
