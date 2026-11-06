package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) GameHistory(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapWatchResults); err != nil {
		httperr.Error(c, err)
		return
	}

	gameIDs, err := h.event.GameHistory(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_ids": gameIDs,
	})
}

func (h Handler) PlayerGameHistory(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	playerID := reqctx.GetPlayerID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapWatchResults); err != nil {
		httperr.Error(c, err)
		return
	}

	gameIDs, err := h.event.PlayerGameHistory(c.Request.Context(), eventID, playerID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_ids": gameIDs,
	})
}
