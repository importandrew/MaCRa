package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func (h Handler) ListMembers(c *gin.Context) {
	eventID := reqctx.GetEventID(c)

	members, err := h.event.ListMembers(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, members)
}
