package event

import (
	"context"
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

type generateTourRequest struct {
	Seating models.SeatingMethod `json:"seating" binding:"required"`
	// Step is an interval seating's step; missing: the number of tables.
	Step *int `json:"step"`
}

// GenerateTour drafts the next tour's seating, replacing an earlier draft
// (POST /event/:id/tours, {"seating": "interval", "step": 2}; the step only
// for "interval"; "fixed" takes the tour from the seating prepared for all
// tours). It needs
// models.CapRunTournament.
func (h Handler) GenerateTour(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var req generateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	tour, err := h.event.GenerateTour(c.Request.Context(), eventID, models.Seating{Method: req.Seating, Step: req.Step})
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tour)
}

// GenerateSeatingPlan prepares the seating for all the remaining tours,
// replacing an earlier plan (POST /event/:id/tours/plan). It needs
// models.CapRunTournament.
func (h Handler) GenerateSeatingPlan(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	info, err := h.event.GenerateSeatingPlan(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// maxSeatingFile bounds an imported seating file: 172 players for 15 tours
// take about 20 KB.
const maxSeatingFile = 256 << 10

type importSeatingPlanRequest struct {
	// CSV is the file's text: "tour,table,east,south,west,north" per table.
	CSV string `json:"csv" binding:"required"`
}

// ImportSeatingPlan stores a seating for several tours from a CSV file,
// replacing an earlier plan (PUT /event/:id/tours/plan, {"csv": "…"}). It
// needs models.CapRunTournament.
func (h Handler) ImportSeatingPlan(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSeatingFile)
	var req importSeatingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	info, err := h.event.ImportSeatingPlan(c.Request.Context(), eventID, req.CSV)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// DeleteSeatingPlan deletes the seating prepared for all tours
// (DELETE /event/:id/tours/plan). It needs models.CapRunTournament.
func (h Handler) DeleteSeatingPlan(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.DeleteSeatingPlan(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelTour cancels the drafted or published seating of the tour not started
// yet (DELETE /event/:id/tours/current). It needs models.CapRunTournament.
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

// PublishTour publishes the drafted seating to the players
// (POST /event/:id/tours/current/publish).
func (h Handler) PublishTour(c *gin.Context) {
	h.tourAction(c, h.event.PublishTour)
}

// StartTour starts the published tour (POST /event/:id/tours/current/start).
func (h Handler) StartTour(c *gin.Context) {
	h.tourAction(c, h.event.StartTour)
}

// FinishTour finishes the live tour once all its games have ended
// (POST /event/:id/tours/current/finish).
func (h Handler) FinishTour(c *gin.Context) {
	h.tourAction(c, h.event.FinishTour)
}

// AbortTour aborts the live tour (POST /event/:id/tours/current/abort).
func (h Handler) AbortTour(c *gin.Context) {
	h.tourAction(c, h.event.AbortTour)
}

// tourAction runs a referee's action on the current tour. It needs
// models.CapRunTournament.
func (h Handler) tourAction(c *gin.Context, action func(ctx context.Context, eventID int) error) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := action(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

type addTimeRequest struct {
	Minutes int `json:"minutes" binding:"required"`
}

// AddTime adds minutes to a table of the live tour
// (POST /event/:id/game/:gameId/time, {"minutes": 5}). It needs
// models.CapRunTournament.
func (h Handler) AddTime(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var req addTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.AddTime(c.Request.Context(), eventID, reqctx.GetGameID(c), req.Minutes); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
