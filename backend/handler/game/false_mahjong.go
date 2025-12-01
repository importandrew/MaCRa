package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type FalseMahjongRequest struct {
	LoserID        int  `json:"loser"`
	Value          int  `json:"value"`
	ForceNextRound bool `json:"force_next_round"`
}

func (h Handler) FalseMahjong(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	var req FalseMahjongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}

	err = h.game.FalseMahjong(c.Request.Context(), id, req.LoserID, req.Value, req.ForceNextRound)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
