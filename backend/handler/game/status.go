package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) Status(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	status, err := h.game.Status(c.Request.Context(), id)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"status": status,
	})
}
