package event

import (
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) HasAccess(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if err := h.event.CanAccess(c.Request.Context(), eventID, userID); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Next()
}

func (h Handler) GameBelongsToEvent(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	gameEventID, err := h.game.EventID(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	if gameEventID != eventID {
		httperr.Error(c, models.ErrGameNotInEvent)
		return
	}

	c.Next()
}
