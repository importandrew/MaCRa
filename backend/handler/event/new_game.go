package event

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type NewGameRequest struct {
	Players []int `json:"players" binding:"required,len=4"`
}

func (h Handler) NewGame(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req NewGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, 0, models.CapCreateGame); err != nil {
		httperr.Error(c, err)
		return
	}

	if !slices.Contains(req.Players, callerID) {
		httperr.Error(c, models.ErrUserNotInPlayers)
		return
	}

	gameID, err := h.event.NewGame(c.Request.Context(), eventID, req.Players)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": gameID,
	})
}
