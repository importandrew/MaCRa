package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
)

type MainInfo struct {
	HasLast    bool `json:"has_last"`
	HasCurrent bool `json:"has_current"`
}

func (h Handler) MainInfo(c *gin.Context) {
	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	last, _ := h.game.LastGameResult(c.Request.Context(), userID)
	now, _ := h.game.Now(c.Request.Context(), userID)

	c.JSON(http.StatusOK, MainInfo{
		HasLast:    last.ID > 0,
		HasCurrent: now.ID > 0,
	})
}
