package event

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// UpdateJoining sets whether the event takes requests to join and its
// player limit (PUT /event/:id/joining). It needs models.CapManageMembers.
func (h Handler) UpdateJoining(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	var settings models.JoiningSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.UpdateJoining(c.Request.Context(), eventID, settings); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// JoinRequests lists the event's pending requests to join, oldest first
// (GET /event/:id/requests). It needs models.CapManageMembers.
func (h Handler) JoinRequests(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	requests, err := h.event.JoinRequests(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, requests)
}

// AcceptJoinRequest makes the requesting user a player
// (POST /event/:id/requests/:playerId/accept). It needs
// models.CapManageMembers.
func (h Handler) AcceptJoinRequest(c *gin.Context) {
	h.answerJoinRequest(c, h.event.AcceptJoinRequest)
}

// DeclineJoinRequest declines the user's request to join
// (POST /event/:id/requests/:playerId/decline). It needs
// models.CapManageMembers.
func (h Handler) DeclineJoinRequest(c *gin.Context) {
	h.answerJoinRequest(c, h.event.DeclineJoinRequest)
}

func (h Handler) answerJoinRequest(c *gin.Context, answer func(c context.Context, eventID, userID int) error) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := answer(c.Request.Context(), eventID, reqctx.GetPlayerID(c)); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// RequestToJoin sends the signed-in user's request to join the event
// (POST /event/:id/request).
func (h Handler) RequestToJoin(c *gin.Context) {
	if err := h.event.RequestToJoin(c.Request.Context(), reqctx.GetEventID(c), reqctx.GetUserID(c)); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// WithdrawJoinRequest deletes the signed-in user's pending request
// (DELETE /event/:id/request).
func (h Handler) WithdrawJoinRequest(c *gin.Context) {
	if err := h.event.WithdrawJoinRequest(c.Request.Context(), reqctx.GetEventID(c), reqctx.GetUserID(c)); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
