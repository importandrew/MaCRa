package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) MainInfo(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if userID == reqctx.ZeroUserID {
		c.JSON(http.StatusOK, models.MainGameInfo{
			HasLast:    false,
			HasCurrent: false,
		})
		return
	}

	mainInfo, err := h.event.MainGameInfo(c.Request.Context(), eventID, userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, mainInfo)
}
