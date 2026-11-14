package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// PenaltyRequest is a penalty or a false mahjong, at the table or from the
// game's view.
type PenaltyRequest struct {
	LoserID        int  `json:"loser"`
	Value          int  `json:"value"`
	ForceNextRound bool `json:"force_next_round"`
	// LogSize is the number of round log entries the table (or the view)
	// showed. A pointer, since 0 is valid but a missing one isn't.
	LogSize *int `json:"log_size" binding:"required"`
}

// Penalty enters a penalty at the caller's table (POST /event/:id/penalty).
func (h Handler) Penalty(c *gin.Context) {
	h.penaltyAtTable(c, models.Penalty)
}

// FalseMahjong enters a false mahjong at the caller's table
// (POST /event/:id/false-mahjong).
func (h Handler) FalseMahjong(c *gin.Context) {
	h.penaltyAtTable(c, models.FalseMahjong)
}

// GamePenalty enters a penalty into the game, as a referee does from the
// game's view (POST /event/:id/game/:gameId/penalty).
func (h Handler) GamePenalty(c *gin.Context) {
	h.penaltyInGame(c, models.Penalty)
}

// GameFalseMahjong enters a false mahjong into the game, as GamePenalty
// (POST /event/:id/game/:gameId/false-mahjong).
func (h Handler) GameFalseMahjong(c *gin.Context) {
	h.penaltyInGame(c, models.FalseMahjong)
}

// penaltyAtTable enters the penalty into the caller's live game in the event
// and answers with their table. The caller needs models.CapEnterPenalty for
// that game: the event's settings may leave penalties to referees.
func (h Handler) penaltyAtTable(c *gin.Context, penaltyType models.Outcome) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req PenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	gameID, err := h.game.CurrentGameID(c.Request.Context(), callerID, eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, gameID, models.CapEnterPenalty); err != nil {
		httperr.Error(c, err)
		return
	}

	now, err := h.game.Penalty(c.Request.Context(), penaltyType, callerID, req.LoserID, req.Value, req.ForceNextRound, eventID, *req.LogSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, now)
}

// penaltyInGame enters the penalty into the game of the URL and answers with
// its live view. It needs models.CapEnterPenalty: a referee in any live game
// of the event, or a seated player as the event's settings allow.
func (h Handler) penaltyInGame(c *gin.Context, penaltyType models.Outcome) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	var req PenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, gameID, models.CapEnterPenalty); err != nil {
		httperr.Error(c, err)
		return
	}

	view, err := h.game.PenaltyInGame(c.Request.Context(), penaltyType, gameID, req.LoserID, req.Value, req.ForceNextRound, *req.LogSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, view)
}
