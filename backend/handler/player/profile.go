package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

func (h Handler) GetProfile(c *gin.Context) {
	userID := reqctx.GetUserID(c)

	player, err := h.player.Get(c.Request.Context(), userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h Handler) UpdateProfile(c *gin.Context) {
	userID := reqctx.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.player.Update(c.Request.Context(), userID, req.Name); err != nil {
		httperr.Error(c, err)
		return
	}

	// Return updated player data
	player, err := h.player.Get(c.Request.Context(), userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, player)
}
