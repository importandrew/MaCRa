package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) FinalResult(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	result, err := h.game.GameResult(c.Request.Context(), id)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
