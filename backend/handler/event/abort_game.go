package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) Abort(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	// AbortGame checks the capability itself: a seated player may abort only
	// while the game has no rounds, which it checks again in the abort.
	if err := h.event.AbortGame(c.Request.Context(), eventID, callerID, gameID); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
