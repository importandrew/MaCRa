package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type EventService interface {
	Create(ctx context.Context, params models.EventCreateParams) (int, error)
	Update(ctx context.Context, params models.UpdateEventParams) error
	Finish(ctx context.Context, eventID int) error
	Get(ctx context.Context, eventID int) (models.Event, error)
	ListVisible(ctx context.Context, userID int) ([]models.EventListItem, error)
	GetMe(ctx context.Context, eventID, userID int) (models.EventMe, error)
	CanAccess(ctx context.Context, eventID, userID int) error
	HasCapability(ctx context.Context, eventID, userID, gameID int, cap models.Capability) error
	AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error
	RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error
	ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error)
	GameHistory(ctx context.Context, eventID int) ([]int, error)
	PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error)
	ListGamesByStatus(ctx context.Context, eventID int, status models.GameStatus, limit, offset int) (models.EventGames, error)
	Ratings(ctx context.Context, eventID int) ([]models.EventRatingEntry, error)

	NewGame(ctx context.Context, eventID int, playerIDs []int) (int, error)
	MainGameInfo(ctx context.Context, eventID, userID int) (models.MainGameInfo, error)
	AbortGame(ctx context.Context, gameID int) error
	UndoLastRound(ctx context.Context, gameID int) error
}

type Handler struct {
	event EventService
	game  GameService
}

func NewHandler(event EventService, game GameService) Handler {
	return Handler{event: event, game: game}
}
