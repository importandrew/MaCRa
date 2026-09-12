package event

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type GameService interface {
	Now(ctx context.Context, playerID, eventID int) (models.GameNow, error)
	LastGameResult(ctx context.Context, playerID, eventID int) (models.GameResult, error)
	EnterRound(ctx context.Context, playerID int, outcome models.Outcome, value, winnerID, loserID, eventID int) (models.GameNow, error)
	Penalty(ctx context.Context, penaltyType models.Outcome, playerID, penalizedID, value int, forceNextRound bool, eventID int) (models.GameNow, error)
	Finish(ctx context.Context, gameID int) error
	ViewLive(ctx context.Context, gameID int) (models.GameView, error)
	GameResult(ctx context.Context, gameID int) (models.GameResult, error)
	Players(ctx context.Context, gameID int) (models.GamePlayers, error)
	Logs(ctx context.Context, gameID int) (models.GameLog, error)
	HUD(ctx context.Context, gameID int, linger int) (models.GameHUD, error)
	Status(ctx context.Context, gameID int) (models.GameStatus, error)
	EventID(ctx context.Context, gameID int) (int, error)
}

func (h Handler) GameViewLive(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, gameID, models.CapWatchLive); err != nil {
		httperr.Error(c, err)
		return
	}

	view, err := h.game.ViewLive(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, view)
}

func (h Handler) GameResult(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, gameID, models.CapWatchResults); err != nil {
		httperr.Error(c, err)
		return
	}

	result, err := h.game.GameResult(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, gameID, models.CapUndoRound); err != nil {
		result.CanUndoLastRound = false
	}

	c.JSON(http.StatusOK, result)
}

func (h Handler) GamePlayers(c *gin.Context) {
	gameID := reqctx.GetGameID(c)

	players, err := h.game.Players(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, players)
}

func (h Handler) GameLogs(c *gin.Context) {
	gameID := reqctx.GetGameID(c)

	logs, err := h.game.Logs(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h Handler) GameHUD(c *gin.Context) {
	var req struct {
		linger int `form:"linger"`
	}

	if err := c.BindQuery(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, gameID, models.CapWatchLive); err != nil {
		httperr.Error(c, err)
		return
	}

	hud, err := h.game.HUD(c.Request.Context(), gameID, req.linger)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, hud)
}

func (h Handler) GameStatus(c *gin.Context) {
	gameID := reqctx.GetGameID(c)

	status, err := h.game.Status(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     gameID,
		"status": status,
	})
}
