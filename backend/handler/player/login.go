package player

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const cookieTTL = 400 * 24 * 60 * 60 // 400 days in sec (max)

func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	player, err := h.player.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	token, err := authmw.GenerateJWTToken(player.ID, h.jwt)
	if err != nil {
		httperr.Error(c, fmt.Errorf("error generating token: %w", err))
	}

	h.writeAuthCookie(c, token, cookieTTL)

	c.JSON(http.StatusOK, LoginResponse{
		ID:   player.ID,
		Name: player.Name,
	})
}
