package player

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/authcookie"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/handler/token"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the signed-in user. The session token travels only in the
// HttpOnly cookie, never in the body, so page scripts can't read it.
type LoginResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	authcookie.Set(c, tkn)
	// The /auth routes run without SetUserID; set it here so the request log
	// shows who signed in.
	c.Set(reqctx.UserIDKey, player.ID)

	c.JSON(http.StatusOK, LoginResponse{
		ID:   player.ID,
		Name: player.Name,
	})
}
