package player

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
)

var (
	errGeneratingToken = errors.New("error generating token")
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const cookieTTL = 90 * 24 * 60 * 60 // 3 months in sec

func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	player, err := h.player.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	token, err := authmw.GenerateJWTToken(player.ID, h.jwt)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, errGeneratingToken)
	}

	if h.domain != "" {
		c.SetCookie(authmw.TokenKey, token, cookieTTL, "/", h.domain, true, true)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(authmw.TokenKey, token, cookieTTL, "/", models.DefaultDomain, true, true)
	}

	c.JSON(http.StatusOK, LoginResponse{
		ID:   player.ID,
		Name: player.Name,
	})
}
