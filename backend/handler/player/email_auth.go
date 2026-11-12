package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
)

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

func (h Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.player.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

func (h Handler) ResendVerification(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.player.ResendVerification(c.Request.Context(), req.Email); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.player.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.InvalidRequest(c, err)
		return
	}

	if err := h.player.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
