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
	if err := c.BindQuery(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
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
