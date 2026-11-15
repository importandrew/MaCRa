package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// Participants lists a tournament's players with their numbers
// (GET /event/:id/participants). It needs models.CapViewTournament.
func (h Handler) Participants(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapViewTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	participants, err := h.event.Participants(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, participants)
}

// Draw gives the tournament's players new random numbers
// (POST /event/:id/draw). It needs models.CapRunTournament.
func (h Handler) Draw(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.Draw(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Withdraw marks a participant as withdrawn
// (POST /event/:id/participants/:playerId/withdraw). It needs
// models.CapRunTournament.
func (h Handler) Withdraw(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.Withdraw(c.Request.Context(), eventID, reqctx.GetPlayerID(c)); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
