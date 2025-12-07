package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type ListNameLimit struct {
	Name  string `form:"name"`
	Limit int    `form:"limit"`
}

func (h Handler) List(c *gin.Context) {
	var req ListNameLimit
	if err := c.BindQuery(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	players, err := h.player.ListByNameLimit(c.Request.Context(), req.Name, req.Limit)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, players)
}
