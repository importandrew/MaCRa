package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) GetMe(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	me, err := h.event.GetMe(c.Request.Context(), eventID, userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, me)
}

// GameMe returns the caller's capabilities for one game: what the game's
// screens may offer (enter rounds, undo, abort, finish).
func (h Handler) GameMe(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	me, err := h.event.MyGameCapabilities(c.Request.Context(), eventID, userID, gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, me)
}
