package game

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
)

type NewGameRequest struct {
	Players []int `json:"players"`
}

func (h Handler) New(c *gin.Context) {
	var req NewGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !slices.Contains(req.Players, userID) {
		httperr.Error(c, models.ErrUserNotInPlayers)
		return
	}

	_, err = h.game.NewGame(c.Request.Context(), req.Players)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
