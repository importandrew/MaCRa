package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) UndoLastRound(c *gin.Context) {
	strID := c.Param("id")
	gameID, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	err = h.game.UndoLastRound(c.Request.Context(), gameID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
