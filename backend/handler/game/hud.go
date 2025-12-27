package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type Linger struct {
	Linger int `form:"linger"`
}

func (h Handler) HUD(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, fmt.Errorf("id should be integer, got: %v", strID))
		return
	}

	var req Linger
	if err := c.BindQuery(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	hud, err := h.game.HUD(c.Request.Context(), id, req.Linger)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, hud)
}
