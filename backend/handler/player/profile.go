package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type ProfileResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

func (h Handler) GetProfile(c *gin.Context) {
	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusUnauthorized, err)
		return
	}

	player, err := h.player.Get(c.Request.Context(), userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, ProfileResponse{
		ID:   player.ID,
		Name: player.Name,
	})
}

func (h Handler) UpdateProfile(c *gin.Context) {
	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusUnauthorized, err)
		return
	}

	var req UpdateProfileRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
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

	c.JSON(http.StatusOK, ProfileResponse{
		ID:   player.ID,
		Name: player.Name,
	})
}
