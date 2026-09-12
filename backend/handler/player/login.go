package player

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/token"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

const cookieTTL = 400 * 24 * 60 * 60 // 400 days in sec (max possible Max-Age)

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

	tkn, err := token.GenerateJWT(player.ID, h.jwt)
	if err != nil {
		httperr.Error(c, fmt.Errorf("error generating token: %w", err))
		return
	}

	h.writeAuthCookie(c, tkn, cookieTTL)

	c.JSON(http.StatusOK, LoginResponse{
		ID:    player.ID,
		Name:  player.Name,
		Token: tkn,
	})
}
