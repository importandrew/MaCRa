package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) Now(c *gin.Context) {
	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	now, err := h.game.Now(c.Request.Context(), userID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, now)
}
