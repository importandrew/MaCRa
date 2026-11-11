package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h Handler) SignUp(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.player.SignUp(c.Request.Context(), req.Email, req.Name, req.Password); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
