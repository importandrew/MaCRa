package game

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
)

type EnterRoundResultRequest struct {
	Outcome  models.Outcome `json:"outcome"`
	Value    int            `json:"score"`
	WinnerID int            `json:"won"`
	LoserID  int            `json:"dealt"`
}

func (h Handler) EnterRoundResult(c *gin.Context) {
	var req EnterRoundResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("bind: %w", err))
		return
	}

	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	err = h.game.EnterRound(c.Request.Context(), userID, req.Outcome, req.Value, req.WinnerID, req.LoserID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
