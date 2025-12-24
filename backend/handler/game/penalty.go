package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type PenaltyRequest struct {
	LoserID        int  `json:"loser"`
	Value          int  `json:"value"`
	ForceNextRound bool `json:"force_next_round"`
}

func (h Handler) Penalty(c *gin.Context) {
	var req PenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	now, err := h.game.Penalty(c.Request.Context(), userID, req.LoserID, req.Value, req.ForceNextRound)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, now)
}
