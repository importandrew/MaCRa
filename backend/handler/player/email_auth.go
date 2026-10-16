package player

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
)

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

func (h Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.player.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type ResendVerificationRequest struct {
	Email  string `json:"email"`
	Origin string `json:"origin"`
}

func (h Handler) ResendVerification(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	origin, err := normalizeOrigin(req.Origin, h.allowOrigins)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.player.ResendVerification(c.Request.Context(), req.Email, origin); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type ForgotPasswordRequest struct {
	Email  string `json:"email"`
	Origin string `json:"origin"`
}

func (h Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	origin, err := normalizeOrigin(req.Origin, h.allowOrigins)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.player.ForgotPassword(c.Request.Context(), req.Email, origin); err != nil {
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
	if err := c.BindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.player.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// normalizeOrigin validates origin and returns scheme://host (no trailing slash/path).
func normalizeOrigin(origin string, allowOrigins []string) (string, error) {
	u, err := url.Parse(origin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", models.ErrInvalidOrigin
	}

	if !slices.Contains([]string{"http", "https"}, u.Scheme) {
		return "", models.ErrInvalidOrigin
	}

	origin = strings.TrimRight(strings.TrimSpace(u.Scheme+"://"+u.Host), "/")

	if !slices.Contains(allowOrigins, origin) {
		return "", models.ErrInvalidOrigin
	}

	return origin, nil
}
