package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
)

// List finds players for the player pickers (GET /api/players?q=): a user ID,
// or part of a name (see Service.Search).
func (h Handler) List(c *gin.Context) {
	var req struct {
		Query string `form:"q"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	players, err := h.player.Search(c.Request.Context(), req.Query)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, players)
}
