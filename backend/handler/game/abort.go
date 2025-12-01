package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) Abort(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	err = h.game.Abort(c.Request.Context(), id)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("game abort: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
