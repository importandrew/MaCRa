package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) UndoLastRound(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, callerID, gameID, models.CapUndoRound); err != nil {
		httperr.Error(c, err)
		return
	}

	err := h.event.UndoLastRound(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
