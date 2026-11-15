package event

import (
	"context"
	"net/http"
	"strconv"

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
	// AvoidRepeats: a Swiss or interval seating avoids repeat meetings.
	AvoidRepeats bool `json:"avoid_repeats"`
}

// GenerateTour drafts the next tour's seating, replacing an earlier draft
// (POST /event/:id/tours, {"seating": "interval", "step": 2,
// "avoid_repeats": true}; the step only for "interval", avoiding repeats
// for "swiss" and "interval"; "fixed" takes the tour's tables from a file). It needs
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

	tour, err := h.event.GenerateTour(c.Request.Context(), eventID, models.Seating{Method: req.Seating, Step: req.Step, AvoidRepeats: req.AvoidRepeats})
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tour)
}

type tourPlanRequest struct {
	Seating models.SeatingMethod `json:"seating" binding:"required"`
	// Step is an interval seating's step; missing: the number of tables.
	Step         *int `json:"step"`
	AvoidRepeats bool `json:"avoid_repeats"`
}

// SetTourPlan sets the seating a tour is to have
// (PUT /event/:id/tours/plan/:number, {"seating": "interval", "step": 2,
// "avoid_repeats": true}). It needs models.CapRunTournament.
func (h Handler) SetTourPlan(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		httperr.InvalidRequest(c, err)
		return
	}
	var req tourPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	plan := models.TourPlan{Number: number, Seating: req.Seating, Step: req.Step, AvoidRepeats: req.AvoidRepeats}
	if err := h.event.SetTourPlan(c.Request.Context(), eventID, plan); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// maxSeatingFile bounds an imported seating file: 172 players for 15 tours
// take about 20 KB.
const maxSeatingFile = 256 << 10

type importSeatingFileRequest struct {
	// CSV is the file's text: "tour,table,east,south,west,north" per table.
	CSV string `json:"csv" binding:"required"`
}

// ImportSeatingFile stores tables for several tours from a CSV file,
// replacing the ones imported before (PUT /event/:id/tours/file,
// {"csv": "…"}). It needs models.CapRunTournament.
func (h Handler) ImportSeatingFile(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSeatingFile)
	var req importSeatingFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	info, err := h.event.ImportSeatingFile(c.Request.Context(), eventID, req.CSV)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// DeleteSeatingFile deletes the tables imported from a file
// (DELETE /event/:id/tours/file). It needs models.CapRunTournament.
func (h Handler) DeleteSeatingFile(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.DeleteSeatingFile(c.Request.Context(), eventID); err != nil {
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

type seatSitterRequest struct {
	// Seat is the seat's player (the seats keep their players' IDs).
	Seat   int `json:"seat" binding:"required"`
	Sitter int `json:"sitter" binding:"required"`
	// Forfeit: the seat's player gets 0 TP and 0 GP for the game.
	Forfeit bool `json:"forfeit"`
	// LogSize is the number of log entries the referee saw.
	LogSize *int `json:"log_size" binding:"required"`
}

// SeatSitter seats a substitute at a seat of a tournament game, or the
// seat's player back, from the next log entry on
// (POST /event/:id/game/:gameId/seat, {"seat": 5, "sitter": 9,
// "forfeit": false, "log_size": 12}). It needs models.CapRunTournament.
func (h Handler) SeatSitter(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var req seatSitterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapRunTournament); err != nil {
		httperr.Error(c, err)
		return
	}

	err := h.game.SeatSitter(c.Request.Context(), models.SeatChange{
		EventID: eventID, GameID: reqctx.GetGameID(c), SeatPlayerID: req.Seat,
		SitterID: req.Sitter, Forfeit: req.Forfeit, SeenLogSize: *req.LogSize,
	})
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
