package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) Get(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	event, err := h.event.Get(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h Handler) ListVisible(c *gin.Context) {
	userID := reqctx.GetUserID(c)

	list, err := h.event.ListVisible(c.Request.Context(), userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}
