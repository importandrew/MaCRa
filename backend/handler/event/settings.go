package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// UpdateSettings sets the event's settings (PUT /event/:id/settings, all of
// them at once). It needs models.CapManageEvent.
func (h Handler) UpdateSettings(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var settings models.EventSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapManageEvent); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.UpdateSettings(c.Request.Context(), eventID, settings); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
