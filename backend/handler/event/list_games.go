package event

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) ListGames(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	status := models.GameStatus(c.Query("status"))
	if status != models.GameLive && status != models.GameEnded {
		httperr.Error(c, models.ErrInvalidGameStatus)
		return
	}

	cap := models.CapWatchResults
	if status == models.GameLive {
		cap = models.CapWatchLive
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, cap); err != nil {
		httperr.Error(c, err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.event.ListGamesByStatus(c.Request.Context(), eventID, status, limit, offset)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
