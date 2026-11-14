package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// PlayerStats returns a player's statistics in the event
// (GET /event/:id/player/:playerId/stats). A player always sees their own;
// access to the event itself is checked by HasAccess. Other players'
// statistics follow the rating: models.CapViewRatings.
func (h Handler) PlayerStats(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	playerID := reqctx.GetPlayerID(c)

	if playerID != userID {
		if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapViewRatings); err != nil {
			httperr.Error(c, err)
			return
		}
	}

	stats, err := h.event.PlayerStats(c.Request.Context(), eventID, playerID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
