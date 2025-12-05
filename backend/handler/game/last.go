package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

func (h Handler) Last(c *gin.Context) {
	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	result, err := h.game.LastGameResult(c.Request.Context(), userID)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
