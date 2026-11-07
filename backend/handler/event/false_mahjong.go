package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type FalseMahjongRequest struct {
	LoserID        int  `json:"loser"`
	Value          int  `json:"value"`
	ForceNextRound bool `json:"force_next_round"`
	// LogSize is the number of round log entries the table showed
	// (GameNow.log_size). A pointer, since 0 is valid but a missing one isn't.
	LogSize *int `json:"log_size" binding:"required"`
}

func (h Handler) FalseMahjong(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req FalseMahjongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	now, err := h.game.Penalty(c.Request.Context(), models.FalseMahjong, callerID, req.LoserID, req.Value, req.ForceNextRound, eventID, *req.LogSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, now)
}
