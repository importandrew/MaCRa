package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// UpdateTournament sets a tournament's settings (PUT /event/:id/tournament,
// all of them at once). It needs models.CapManageEvent.
func (h Handler) UpdateTournament(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var settings models.TournamentSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageEvent); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.UpdateTournament(c.Request.Context(), eventID, settings); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
