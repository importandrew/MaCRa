package event

import (
	"errors"
	"net/http"
	"strconv"

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

	// log_size is the number of round log entries the referee saw; the undo
	// is refused if the log has changed since (see repo UndoLastRound).
	logSize, err := strconv.Atoi(c.Query("log_size"))
	if err != nil || logSize < 0 {
		httperr.InvalidRequest(c, errors.New("log_size query parameter is required"))
		return
	}

	err = h.event.UndoLastRound(c.Request.Context(), gameID, logSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
