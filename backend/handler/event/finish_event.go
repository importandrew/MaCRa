package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) FinishEvent(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapManageEvent); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.Finish(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
