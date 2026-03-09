package player

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) GameHistory(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	gameIDs, err := h.player.GameHistory(c.Request.Context(), id)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PlayerHistory{
		ID:      id,
		GameIDs: gameIDs,
	})
}
