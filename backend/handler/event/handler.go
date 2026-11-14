package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type EventService interface {
	Create(ctx context.Context, params models.EventCreateParams) (int, error)
	Update(ctx context.Context, params models.UpdateEventParams) error
	UpdateRules(ctx context.Context, eventID int, rules models.GameRules) error
	UpdateSettings(ctx context.Context, eventID int, settings models.EventSettings) error
	Finish(ctx context.Context, eventID int) error
	Details(ctx context.Context, eventID int) (models.EventDetails, error)
	ListVisible(ctx context.Context, userID int) ([]models.EventListItem, error)
	GetMe(ctx context.Context, eventID, userID int) (models.EventMe, error)
	MyGameCapabilities(ctx context.Context, eventID, callerID, gameID int) (models.GameMe, error)
	CanAccess(ctx context.Context, eventID, userID int) error
	HasCapability(ctx context.Context, eventID, userID, gameID int, capability models.Capability) error
	AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error
	RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error
	SetMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error
	InviteToken(ctx context.Context, eventID int) (*string, error)
	NewInviteToken(ctx context.Context, eventID int) (string, error)
	TurnOffInvite(ctx context.Context, eventID int) error
	Invite(ctx context.Context, token string, userID int) (models.EventInvite, error)
	JoinByInvite(ctx context.Context, token string, userID int) (int, error)
	UpdateJoining(ctx context.Context, eventID int, settings models.JoiningSettings) error
	UpdateTournament(ctx context.Context, eventID int, settings models.TournamentSettings) error
	RequestToJoin(ctx context.Context, eventID, userID int) error
	WithdrawJoinRequest(ctx context.Context, eventID, userID int) error
	JoinRequests(ctx context.Context, eventID int) ([]models.JoinRequest, error)
	AcceptJoinRequest(ctx context.Context, eventID, userID int) error
	DeclineJoinRequest(ctx context.Context, eventID, userID int) error
	ListMembers(ctx context.Context, eventID, callerID int) ([]models.EventMember, error)
	GameHistory(ctx context.Context, eventID int) ([]int, error)
	PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error)
	PlayerStats(ctx context.Context, eventID, playerID int) (models.PlayerStats, error)
	ListGamesByStatus(ctx context.Context, eventID int, status models.GameStatus, limit, offset int) (models.EventGames, error)
	Ratings(ctx context.Context, eventID, callerID int) (models.EventRatings, error)

	NewGame(ctx context.Context, eventID int, playerIDs []int) (int, error)
	MainGameInfo(ctx context.Context, eventID, userID int) (models.MainGameInfo, error)
	AbortGame(ctx context.Context, eventID, callerID, gameID int) error
	UndoLastRound(ctx context.Context, eventID, callerID, gameID, logSize int) error
}

type Handler struct {
	event EventService
	game  GameService
}

func NewHandler(event EventService, game GameService) Handler {
	return Handler{event: event, game: game}
}
