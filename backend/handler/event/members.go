package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) RemoveMember(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req struct {
		IDs []int `form:"ids"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.RemoveMembers(c.Request.Context(), eventID, req.IDs); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type setMemberRolesRequest struct {
	Roles []models.EventRole `json:"roles" binding:"required"`
}

// SetMemberRoles sets a member's roles, all of them at once
// (PUT /event/:id/members/:playerId). It needs models.CapManageMembers.
func (h Handler) SetMemberRoles(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	userID := reqctx.GetPlayerID(c)

	var req setMemberRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.SetMemberRoles(c.Request.Context(), eventID, userID, req.Roles); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
