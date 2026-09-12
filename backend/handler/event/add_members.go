package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type addMembersRequest struct {
	UserID int                `json:"user_id" binding:"required"`
	Roles  []models.EventRole `json:"roles" binding:"required"`
}

func (h Handler) AddMembers(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req addMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.AddMemberRoles(c.Request.Context(), eventID, req.UserID, req.Roles); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
