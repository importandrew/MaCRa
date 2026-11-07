package event

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type EnterRoundResultRequest struct {
	Outcome  models.Outcome `json:"outcome"`
	Value    int            `json:"score"`
	WinnerID int            `json:"won"`
	LoserID  int            `json:"dealt"`
	// LogSize is the number of round log entries the table showed
	// (GameNow.log_size). A pointer, since 0 is valid but a missing one isn't.
	LogSize *int `json:"log_size" binding:"required"`
}

func (h Handler) EnterRoundResult(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req EnterRoundResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}

	now, err := h.game.EnterRound(c.Request.Context(), callerID, req.Outcome, req.Value, req.WinnerID, req.LoserID, eventID, *req.LogSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, now)
}
