package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) Players(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	players, err := h.game.Players(c.Request.Context(), id)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, players)
}
