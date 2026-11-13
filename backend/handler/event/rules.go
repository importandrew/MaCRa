package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// UpdateRules sets the event's game rules (PUT /event/:id/rules, all of them
// at once). It needs models.CapManageEvent.
func (h Handler) UpdateRules(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var rules models.GameRules
	if err := c.ShouldBindJSON(&rules); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapManageEvent); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.UpdateRules(c.Request.Context(), eventID, rules); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
