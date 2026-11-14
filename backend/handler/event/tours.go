package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// Tours returns the tournament's current tour with its tables and what the
// next tour would be (GET /event/:id/tours). It needs
// models.CapViewTournament.
func (h Handler) Tours(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapViewTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	state, err := h.event.Tours(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, state)
}

type publishTourRequest struct {
	Seating models.SeatingMethod `json:"seating" binding:"required"`
}

// PublishTour publishes the next tour's seating (POST /event/:id/tours,
// {"seating": "random"}). It needs models.CapRunTournament.
func (h Handler) PublishTour(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var req publishTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	tour, err := h.event.PublishTour(c.Request.Context(), eventID, req.Seating)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tour)
}

// CancelTour cancels the published seating of the tour not started yet
// (DELETE /event/:id/tours/current). It needs models.CapRunTournament.
func (h Handler) CancelTour(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.CancelTour(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
