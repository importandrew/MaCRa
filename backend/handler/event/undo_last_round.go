package event

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) UndoLastRound(c *gin.Context) {
	callerID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)
	gameID := reqctx.GetGameID(c)

	// log_size is the number of round log entries the caller saw; the undo
	// is refused if the log has changed since (see repo UndoLastRound).
	logSize, err := strconv.Atoi(c.Query("log_size"))
	if err != nil || logSize < 0 {
		httperr.InvalidRequest(c, errors.New("log_size query parameter is required"))
		return
	}

	// UndoLastRound checks the capability itself: a seated player's undo must
	// not revive a game that has ended since, which it checks in the undo.
	err = h.event.UndoLastRound(c.Request.Context(), eventID, callerID, gameID, logSize)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
