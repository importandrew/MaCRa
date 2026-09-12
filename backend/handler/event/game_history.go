package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) GameHistory(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

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
	eventID := reqctx.GetEventID(c)
	playerID := reqctx.GetPlayerID(c)

	gameIDs, err := h.event.PlayerGameHistory(c.Request.Context(), eventID, playerID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_ids": gameIDs,
	})
}
