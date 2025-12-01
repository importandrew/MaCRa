package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) Logout(c *gin.Context) {
	if h.domain != "" {
		c.SetCookie(authmw.TokenKey, "", -1, "/", h.domain, true, true)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(authmw.TokenKey, "", -1, "/", models.DefaultDomain, true, true)
	}

	c.Status(http.StatusOK)
}
